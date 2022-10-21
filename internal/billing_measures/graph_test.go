package billing_measures

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Increment struct {
	t *Test
}

func (l Increment) Execute(ctx context.Context) error {
	l.t.Value++

	return nil
}

func (l Increment) ID() string {
	return "increment"
}

type Test struct {
	Value int
}

func Test_Unit_Domain_BillingMeasures_Graph_Apply(t *testing.T) {

	type wantNode struct {
		status NodeStatus
		done   bool
	}

	type Want struct {
		t     Test
		nodes map[NodeKey]wantNode
	}

	type TestCase struct {
		input func() (*Graph, *Test)
		want  Want
	}

	tests := map[string]TestCase{
		"simple": {
			input: func() (*Graph, *Test) {
				test := &Test{
					Value: 1,
				}

				nodeTrue := &Node{
					Id:           "node_0",
					Precondition: Simple{true},
					Algorithms:   []Algorithm{Increment{test}},
				}
				nodeFalse := &Node{
					Id:           "node_1",
					Precondition: Simple{false},
					Algorithms:   []Algorithm{Increment{test}},
				}

				return NewGraph().AddVector(nodeTrue, nil).AddVector(nodeFalse, nil), test
			},
			want: Want{
				t: Test{
					Value: 2,
				},
				nodes: map[NodeKey]wantNode{
					"node_0": {
						status: NodeStatusSuccess,
						done:   true,
					},
					"node_1": {
						status: "",
						done:   false,
					},
				},
			},
		},
		"two levels": {
			input: func() (*Graph, *Test) {
				test := &Test{
					Value: 1,
				}

				node0 := &Node{
					Id:           "node_0",
					Precondition: Simple{true},
					Algorithms:   []Algorithm{Increment{test}},
				}
				node1 := &Node{
					Id:           "node_1",
					Precondition: Simple{false},
					Algorithms:   []Algorithm{Increment{test}},
				}
				node3 := &Node{
					Id:           "node_3",
					Precondition: Simple{true},
					Algorithms:   []Algorithm{Increment{test}},
				}
				return NewGraph().AddVector(node0, nil).AddVector(node1, nil).AddVector(node0, node3), test
			},
			want: Want{
				t: Test{
					Value: 2,
				},
				nodes: map[NodeKey]wantNode{
					"node_0": {
						status: "",
						done:   true,
					},
					"node_1": {
						status: "",
						done:   false,
					},
					"node_3": {
						status: NodeStatusSuccess,
						done:   true,
					},
				},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			g, testW := testCase.input()
			g.Execute(context.Background())
			assert.Equal(t, testCase.want.t.Value, testW.Value)
			assert.NotNil(t, g.StartedAt)
			assert.NotNil(t, g.FinishedAt)
			for i, w := range testCase.want.nodes {
				n := g.Dict[i]
				assert.Equal(t, w.status, n.Status)
				assert.Equal(t, w.done, n.Done)
				assert.NotNil(t, n.StartedAt)
				assert.NotNil(t, n.FinishedAt)
			}
		})
	}
}
