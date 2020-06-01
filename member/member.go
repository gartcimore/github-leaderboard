package members
import (
  "fmt"
)

type member struct {

  username string
  avatar string
  profile string
  contributions int
  repositories int
  issues int
}

func New(username string, avatar string, profile string, contributions int,
  repositories int, issues int) member {
    m := member {username, avatar, profile, contributions, repositories, issues}
    return m
}

func (m member) LeavesRemaining() {
    fmt.Printf("%s has %d contributions, %d repositories and %d issues\n", m.username, m.contributions, m.repositories, m.issues)
}
