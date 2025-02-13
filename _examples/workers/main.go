package main

import (
	"fmt"
	"log"

	"github.com/littledevels/go-diagrams/diagram"
	"github.com/littledevels/go-diagrams/nodes/gcp"
)

func main() {
	workerCount := 5

	d, err := diagram.New(
		diagram.Label("Stack"),
		diagram.Filename("workers"),
		diagram.Direction("TB"),
		diagram.WithAttribute("compound", "true"),
	)

	if err != nil {
		log.Fatal(err)
	}

	lb := gcp.Network.LoadBalancing(diagram.NodeLabel("nlb"))
	d.Add(lb)

	db := gcp.Database.Sql(diagram.NodeLabel("db"))
	d.Add(db)

	workers := make([]*diagram.Node, workerCount)

	for i := 0; i < workerCount; i++ {
		label := fmt.Sprintf("worker %d", i+1)
		workers[i] = gcp.Compute.ComputeEngine(diagram.NodeLabel(label))
	}

	backends := make([]*diagram.Node, workerCount)
	for i := 0; i < workerCount; i++ {
		label := fmt.Sprintf("backend %d", i+1)
		backends[i] = gcp.Compute.ComputeEngine(diagram.NodeLabel(label))
	}

	group := diagram.NewGroup("workers")
	group.Label("Workers")
	group.Add(workers...)
	group.ConnectAllTo(db.ID(), diagram.WithEdgeAttribute("ltail", group.ID()))
	group.ConnectAllFrom(lb.ID(), diagram.WithEdgeAttribute("lhead", group.ID()))
	d.Group(group)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
