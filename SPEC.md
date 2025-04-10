# Project Structure

We seem to be a little confused about how this project is structured. Look at this:

```bash
➜  vicohome git:(main) ✗ pwd
/Users/josh/Code/skunkworks/vicohome/vicohome
➜  vicohome git:(main) ✗ tree .
.
├── CLAUDE.md
├── README.md
├── cmd
│   ├── auth.go
│   ├── devices
│   │   ├── get.go
│   │   ├── list.go
│   │   └── root.go
│   ├── events.go
│   └── root.go
├── go.mod
├── go.sum
├── main.go
├── pkg
│   └── auth
│       └── auth.go
├── testing
│   └── test.sh
└── vico-cli
```

We have a `./cmd/devices` module, which seems like a great abstraction. Its like a little sub-cli that we hoist up to the larger one and give some shared code like auth, I love it.

But then we also have these other files, and it starts to make less sense. Lets make sure that our Cobra CLI is following best practices and is implemented consistently. I really like how we have the `./cmd/devices` module set up- lets do that too for events and auth maybe? Make a folder for `./cmd/events` and `./cmd/auth`, setting them up like you have done for `./cmd/devices`. update other code that references this, and lets try clean this up and make it pretty.

# CLI Interface

## Events

Observe how the interface for Events appears. What we must do is make this interface more like the interface for Devices. 

### List Events

#### List All Events (default 24 hours)
./vico-cli events

* Must become:
  - ./vico-cli events list
  - ./vico-cli events list --format table (default)
  - ./vico-cli events list --format json

#### List Events for Last N Hours
./vico-cli events --hours 1

### Get Single Event
./vico-cli event 018594221744243886k4jua3TyFQq

* Must become:
  - ./vico-cli events get 018594221744243886k4jua3TyFQq
  - ./vico-cli events get 018594221744243886k4jua3TyFQq --format table (default)
  - ./vico-cli events get 018594221744243886k4jua3TyFQq --format json

## Devices

This is the model for what we want to have. Please make the Events interface more like this

### List Devices

#### List All Devices (with format options)
./vico-cli devices list
./vico-cli devices list --format table (default)
./vico-cli devices list --format json

#### Get Single Device (with format options)
./vico-cli devices get 854396ddc826ed6e3e4263fa067ee288
./vico-cli devices get 854396ddc826ed6e3e4263fa067ee288 --format table (default)
./vico-cli devices get 854396ddc826ed6e3e4263fa067ee288 --format json