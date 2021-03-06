package input

import (
	"context"
	"time"

	"github.com/Jeffail/benthos/v3/internal/interop"
	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/types"
)

//------------------------------------------------------------------------------

func init() {
	Constructors[TypeResource] = TypeSpec{
		constructor: fromSimpleConstructor(NewResource),
		Summary: `
Resource is an input type that runs a resource input by its name.`,
		Description: `
This input allows you to reference the same configured input resource in multiple places, and can also tidy up large nested configs. For
example, the config:

` + "```yaml" + `
input:
  broker:
    inputs:
      - kafka:
          addresses: [ TODO ]
          topics: [ foo ]
          consumer_group: foogroup
      - gcp_pubsub:
          project: bar
          subscription: baz
` + "```" + `

Could also be expressed as:

` + "```yaml" + `
input:
  broker:
    inputs:
      - resource: foo
      - resource: bar

input_resources:
  - label: foo
    kafka:
      addresses: [ TODO ]
      topics: [ foo ]
      consumer_group: foogroup

  - label: bar
    gcp_pubsub:
      project: bar
      subscription: baz
 ` + "```" + `

You can find out more about resources [in this document.](/docs/configuration/resources)`,
		Categories: []Category{
			CategoryUtility,
		},
	}
}

//------------------------------------------------------------------------------

// Resource is an input that wraps an input resource.
type Resource struct {
	mgr          types.Manager
	name         string
	log          log.Modular
	mErrNotFound metrics.StatCounter
}

// NewResource returns a resource input.
func NewResource(
	conf Config, mgr types.Manager, log log.Modular, stats metrics.Type,
) (Type, error) {
	if err := interop.ProbeInput(context.Background(), mgr, conf.Resource); err != nil {
		return nil, err
	}
	return &Resource{
		mgr:          mgr,
		name:         conf.Resource,
		log:          log,
		mErrNotFound: stats.GetCounter("error_not_found"),
	}, nil
}

//------------------------------------------------------------------------------

// TransactionChan returns a transactions channel for consuming messages from
// this input type.
func (r *Resource) TransactionChan() (tChan <-chan types.Transaction) {
	if err := interop.AccessInput(context.Background(), r.mgr, r.name, func(i types.Input) {
		tChan = i.TransactionChan()
	}); err != nil {
		r.log.Debugf("Failed to obtain input resource '%v': %v", r.name, err)
		r.mErrNotFound.Incr(1)
	}
	return
}

// Connected returns a boolean indicating whether this input is currently
// connected to its target.
func (r *Resource) Connected() (isConnected bool) {
	if err := interop.AccessInput(context.Background(), r.mgr, r.name, func(i types.Input) {
		isConnected = i.Connected()
	}); err != nil {
		r.log.Debugf("Failed to obtain input resource '%v': %v", r.name, err)
		r.mErrNotFound.Incr(1)
	}
	return
}

// CloseAsync shuts down the processor and stops processing requests.
func (r *Resource) CloseAsync() {
}

// WaitForClose blocks until the processor has closed down.
func (r *Resource) WaitForClose(timeout time.Duration) error {
	return nil
}

//------------------------------------------------------------------------------
