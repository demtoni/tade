package manager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/sethvargo/go-password/password"
)

const (
	defaultAddress = "0.0.0.0"
	defaultMethod  = "chacha20-ietf-poly1305"
	defaultBackend = "ssserver"
)

var (
	Hostname    string
	PathToState string
	Secret      string
	Addr        string
)

type Manager struct {
	addr      *net.UDPAddr
	portRange [2]int
	state     map[int]*Server
	mutex     sync.RWMutex
}

func New() (*Manager, error) {
	m := &Manager{}

	return m, m.loadState()
}

type LocalState struct {
	PortRange [2]int     `json:"port_range"`
	State     []*Options `json:"state,omitempty"`
}

func (m *Manager) loadState() error {
	in, err := os.Open(PathToState)
	if err != nil {
		return err
	}

	if m.addr, err = net.ResolveUDPAddr("udp", Addr); err != nil {
		return err
	}

	m.state = make(map[int]*Server, 0)

	local := &LocalState{}

	if err := json.NewDecoder(in).Decode(local); err != nil {
		return err
	}

	m.portRange = local.PortRange

	for i := m.portRange[0]; i < m.portRange[1]; i++ {
		m.state[i] = nil
	}

	for k := range local.State {
		s := &Server{local.State[k], nil}

		if s.opts.Port < 0 || s.opts.Port > 0xffff {
			return fmt.Errorf("%s: %s: bad port number\n", PathToState, s.opts.Name)
		}

		m.state[s.opts.Port] = s
	}

	var wg sync.WaitGroup
	chanErr := make(chan error, 1)

	for _, s := range m.state {
		wg.Add(1)
		go func(s *Server, wg *sync.WaitGroup) {
			if s == nil {
				wg.Done()
				return
			}

			var err error

			if err = s.spawn(); err != nil {
				chanErr <- fmt.Errorf("failed to spawn a server for port %d: %s", s.opts.Port, err)
				wg.Done()
				return
			}

			wg.Done()
		}(s, &wg)
	}

	go func() {
		wg.Wait()
		close(chanErr)
	}()

	for err = range chanErr {
		if err != nil {
			log.Printf("warning: %s", err)
		}
	}

	wg.Wait()

	return nil
}

func (m *Manager) saveState() error {
	log.Println("saving current state")
	out, err := os.Create(PathToState)
	if err != nil {
		return err
	}

	m.mutex.Lock()

	local := &LocalState{PortRange: m.portRange}
	local.State = make([]*Options, 0)
	for _, v := range m.state {
		if v != nil {
			local.State = append(local.State, v.opts)
		}
	}

	if err := json.NewEncoder(out).Encode(local); err != nil {
		return err
	}

	m.mutex.Unlock()

	return nil
}

func (m *Manager) get(name string) *Server {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for k, v := range m.state {
		if v != nil && v.opts.Name == name {
			return m.state[k]
		}
	}
	return nil
}

func (m *Manager) getFreePort() (int, error) {
	m.mutex.RLock()

	port := 0
	for k, v := range m.state {
		if v == nil {
			port = k
			break
		}
	}

	m.mutex.RUnlock()

	if port == 0 {
		return 0, fmt.Errorf("couldn't find free port")
	}

	return port, nil
}

func (m *Manager) add(opts *Options) error {
	switch {
	case opts.Name == "":
		return fmt.Errorf("name is empty")
	case opts.Pass == "":
		return fmt.Errorf("password is empty")
	}
	if opts.Addr == "" {
		opts.Addr = defaultAddress
	}
	if opts.Method == "" {
		opts.Method = defaultMethod
	}
	if opts.Backend == "" {
		opts.Backend = defaultBackend
	}

	var err error

	opts.Port, err = m.getFreePort()
	if err != nil {
		return err
	}

	if m.get(opts.Name) != nil {
		return fmt.Errorf("name is already taken")
	}

	s := &Server{opts, nil}
	if err = s.spawn(); err != nil {
		return err
	}

	m.mutex.Lock()

	m.state[opts.Port] = s

	m.mutex.Unlock()

	return nil
}

func (m *Manager) remove(name string) error {
	s := m.get(name)
	if s == nil {
		return fmt.Errorf("server is not running for %s", name)
	}
	log.Println("stopping server for", name)
	if err := s.kill(); err != nil {
		return err
	}

	m.mutex.Lock()

	m.state[s.opts.Port] = nil

	m.mutex.Unlock()

	return nil
}

type Server struct {
	opts *Options
	cmd  *exec.Cmd
}

type Options struct {
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Pass    string `json:"password"`
	Addr    string `json:"addr"`
	Method  string `json:"method"`
	Backend string `json:"backend"`
	Plugin  string `json:"plugin"`
}

var backends = map[string]string{
	"ssserver": "-s %s:%d -k %s -m %s",
}

var plugins = map[string]string{
	"v2ray": "--plugin v2ray-plugin --plugin-opts server",
	"none":  "",
}

func (s *Server) spawn() error {
	name, err := exec.LookPath(s.opts.Backend)
	if err != nil {
		return fmt.Errorf("couldn't find the location of %s", s.opts.Backend)
	}

	argv := strings.Split(fmt.Sprintf(backends[s.opts.Backend], s.opts.Addr,
		s.opts.Port, s.opts.Pass, s.opts.Method), " ")
	if plugins[s.opts.Plugin] != "" {
		argv = append(argv, plugins[s.opts.Plugin])
	}
	log.Printf("name: %s, argv: %v", name, argv)

	s.cmd = exec.Command(name, argv...)
	if err := s.cmd.Start(); err != nil {
		s.cmd = nil
		return err
	}

	return nil
}

func (s *Server) kill() error {
	err := s.cmd.Process.Kill()
	if err != nil {
		return err
	}

	s.cmd.Wait()

	return nil
}

func (m *Manager) Serve() error {
	http.HandleFunc(fmt.Sprintf("POST /%s/", Secret), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pass, err := password.Generate(16, rand.Intn(7), 0, false, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		opts := &Options{
			Name:   r.PostFormValue("name"),
			Pass:   pass,
			Method: r.PostFormValue("method"),
			Plugin: r.PostFormValue("plugin"),
		}
		log.Println(*opts)
		if err := m.add(opts); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		m.saveState()
	})

	http.HandleFunc(fmt.Sprintf("GET /%s/{name}", Secret), func(w http.ResponseWriter, r *http.Request) {
		s := m.get(r.PathValue("name"))
		if s == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uri := "ss://" + base64.StdEncoding.WithPadding(base64.NoPadding).
			EncodeToString([]byte(fmt.Sprintf("%s:%s@%s:%d", s.opts.Method, s.opts.Pass, Hostname, s.opts.Port)))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"connect_url\":\"%s\"}", uri)))
	})

	http.HandleFunc(fmt.Sprintf("DELETE /%s/{name}", Secret), func(w http.ResponseWriter, r *http.Request) {
		if err := m.remove(r.PathValue("name")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		m.saveState()
	})

	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(Addr, nil)
}
