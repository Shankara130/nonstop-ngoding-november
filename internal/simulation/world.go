package simulation

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type World struct {
	Size      float64
	Agents    map[int]*Agent
	Resources map[int]*Resource
	Zones     map[int]*Zone
	nextID    int
	mu        sync.RWMutex
	ticker    *time.Ticker
	stats     WorldStats
}

type WorldStats struct {
	Tick          int     `json:"tick"`
	AgentCount    int     `json:"agentCount"`
	ResourceCount int     `json:"resourceCount"`
	AvgEnergy     float64 `json:"avgEnergy"`
}

type WorldSnapshot struct {
	Agents    []AgentSnapshot    `json:"agents"`
	Resources []ResourceSnapshot `json:"resources"`
	Zones     []ZoneSnapshot     `json:"zones"`
	Stats     WorldStats         `json:"stats"`
}

func NewWorld(size float64) *World {
	w := &World{
		Size:      size,
		Agents:    make(map[int]*Agent),
		Resources: make(map[int]*Resource),
		Zones:     make(map[int]*Zone),
		nextID:    1,
	}

	// spawn zones
	w.addZone(20, 20, 10, ZoneTypeHome)
	w.addZone(80, 80, 8, ZoneTypeFood)
	w.addZone(50, 20, 6, ZoneTypeWater)
	w.addZone(20, 80, 6, ZoneTypeDanger)

	// spawn initial resources
	for i := 0; i < 15; i++ {
		w.spawnResource()
	}

	return w
}

func (w *World) addZone(x, y, radius float64, zType ZoneType) {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := w.nextID
	w.nextID++
	w.Zones[id] = &Zone{
		ID:     id,
		X:      x,
		Y:      y,
		Radius: radius,
		Type:   zType,
	}
}

func (w *World) spawnResource() {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := w.nextID
	w.nextID++
	w.Resources[id] = &Resource{
		ID:     id,
		X:      rand.Float64() * w.Size,
		Y:      rand.Float64() * w.Size,
		Amount: 50 + rand.Float64()*50,
		Type:   ResourceType(rand.Intn(2)),
	}
}

func (w *World) SpawnAgent(agentType AgentType) {
	w.mu.Lock()
	id := w.nextID
	w.nextID++
	w.mu.Unlock()

	agent := NewAgent(id, w.Size, agentType)

	w.mu.Lock()
	w.Agents[id] = agent
	w.mu.Unlock()

	// every agent runs in its own goroutine
	go agent.Run(w)
}

func (w *World) Run(tickRate time.Duration, onSnapshot func(*WorldSnapshot)) {
	w.ticker = time.NewTicker(tickRate)
	defer w.ticker.Stop()

	resourceSpawnTicker := time.NewTicker(5 * time.Second)
	defer resourceSpawnTicker.Stop()

	for {
		select {
		case <-w.ticker.C:
			w.updateStats()
			snapshot := w.Snapshot()
			onSnapshot(snapshot)
		case <-resourceSpawnTicker.C:
			// spawn new resource if below 20
			w.mu.RLock()
			count := len(w.Resources)
			w.mu.RUnlock()

			if count < 20 {
				w.spawnResource()
			}
		}
	}
}

func (w *World) updateStats() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.stats.Tick++
	w.stats.AgentCount = len(w.Agents)
	w.stats.ResourceCount = len(w.Resources)

	totalEnergy := 0.0
	for _, agent := range w.Agents {
		totalEnergy += agent.Energy
	}
	if len(w.Agents) > 0 {
		w.stats.AvgEnergy = totalEnergy / float64(len(w.Agents))
	}
}

func (w *World) Snapshot() *WorldSnapshot {
	w.mu.RLock()
	defer w.mu.RUnlock()

	agents := make([]AgentSnapshot, 0, len(w.Agents))
	for _, a := range w.Agents {
		agents = append(agents, a.Snapshot())
	}

	resources := make([]ResourceSnapshot, 0, len(w.Resources))
	for _, r := range w.Resources {
		resources = append(resources, r.Snapshot())
	}

	zones := make([]ZoneSnapshot, 0, len(w.Zones))
	for _, z := range w.Zones {
		zones = append(zones, ZoneSnapshot{
			ID:     z.ID,
			X:      z.X,
			Y:      z.Y,
			Radius: z.Radius,
			Type:   z.Type,
		})
	}

	return &WorldSnapshot{
		Agents:    agents,
		Resources: resources,
		Zones:     zones,
		Stats:     w.stats,
	}
}

func (w *World) FindNearestResource(x, y float64, rType ResourceType) *Resource {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var nearest *Resource
	minDist := math.MaxFloat64

	for _, r := range w.Resources {
		if r.Amount <= 0 || r.Type != rType {
			continue
		}
		dist := distance(x, y, r.X, r.Y)
		if dist < minDist {
			minDist = dist
			nearest = r
		}
	}
	return nearest
}

func (w *World) FindNearestZone(x, y float64, zType ZoneType) *Zone {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var nearest *Zone
	minDist := math.MaxFloat64

	for _, z := range w.Zones {
		if z.Type != zType {
			continue
		}
		dist := distance(x, y, z.X, z.Y)
		if dist < minDist {
			minDist = dist
			nearest = z
		}
	}
	return nearest
}

func (w *World) ConsumeResource(resourceID int, amount float64) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	r, exists := w.Resources[resourceID]
	if !exists || r.Amount <= 0 {
		return false
	}

	consumed := math.Min(amount, r.Amount)
	r.Amount -= consumed

	if r.Amount <= 0 {
		delete(w.Resources, resourceID)
	}

	return true
}

func (w *World) RemoveAgent(id int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.Agents, id)
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx, dy := x1-x2, y1-y2
	return math.Sqrt(dx*dx + dy*dy)
}
