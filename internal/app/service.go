package app

// Service aggregates application use-cases.
type Service struct {
    Rooms RoomRepo
    Ids   IdGen
    Bus   Broadcaster
}
