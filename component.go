package imagetemplate

import ()

type Component interface {
	Write(canvas Canvas) error
}
