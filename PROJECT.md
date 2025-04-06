# Project

Extraction of API Client from Application Source Code found in `./app`

API In question:

* api-us.vicohome.io

Based on extensive review of the codebase and its usage of the API, we will begin to catalog these remote endpoints in an `openapi.yml`

## Project Orchestration

We break things down into phases with tasks. Phases can have these states:

* []
* [started]
* [complete]

When we open the project up, we look for the first phase that is marked as [started]. We then review its tasks. Upon completion of all tasks, we mark the phase as [completed]. If there are none marked as [started], we look for the first one marked as []. We update it to [started].


```mermaid
flowchart TD
    A[Start Process] --> B["Any phases marked as [started]?"]
    B -->|Yes| C["Continue from first [started] phase"]
    B -->|No| D["Any phases marked as []?"]
    D -->|Yes| E["Mark first [] phase as [started]"]
    E --> F["Continue from this phase"]
    D -->|No| G["All phases are [complete]"]
    C --> H["Work on phase"]
    F --> H
    H --> I["Mark phase as [complete]"]
    I --> J["More phases remaining?"]
    J -->|Yes| B
    J -->|No| K["End Process"]
    G --> K
```

Likewise, tasks within these phases can have statuses:

* []
* [started]
* [complete]

When we enter a [started] phase, we iterate through the tasks looking for the first one that is marked [started]. Upon completion of the task, we mark it as [completed]. If there are none marked as [started], we look for the first one marked as []. We update it to [started].

```mermaid
flowchart TD
    A[Start Process] --> B["Any tasks marked as [started]?"]
    B -->|Yes| C["Continue from first [started] task"]
    B -->|No| D["Any tasks marked as []?"]
    D -->|Yes| E["Mark first [] task as [started]"]
    E --> F["Continue from this task"]
    D -->|No| G["All tasks are [complete]"]
    C --> H["Work on task"]
    F --> H
    H --> I["Mark task as [complete]"]
    I --> J["More task remaining?"]
    J -->|Yes| B
    J -->|No| K["End Process"]
    G --> K
```

# Phase One [complete] Research API

1. [complete] Create `endpoints.yml` that contains a list of all known endpoints. We create a grouped inventory of endpoints. Each new category has some added metadata about its analysis state. For now, our state is "unverified". Each new endpoint added has metadata accounting for its analysis state. For now, our state is "unverified"

2. [started] Verify endpoints in `endpoints.yml`. Here we look for the first api group that states it is "unverified". We will then search for the first api endpoint under that group that is "unverified". We will do deep analysis of the code found in `/app` for any and all information about how this endpoint is called, what parameters are being sent, what is being obtained as a result of the request, and whether any error handling is occurring. With what we learn about the endpoint during this evaluation, we create a special markdown file inside of the `endpoints/` folder in the form of `${nice-endpoint-name}-discovery.md`. When a group has endpoints who's statuses are all `openapi`, we may mark the group as `openapi`.

3. [] Form `openapi.yml`. We will review the contents of `endpoints.yml` and identify the first one that is marked as "verified". We will find the corresponding research information about this endpoint in `${nice-endpoint-name}-discovery.md`. With this information we will create the appropriate `paths` and `components` entries necessary to properly add this endpoint to our `openapi.yml`. Once we know what we will place in `openapi.yml`, we will add the necessary `paths` and `components` updates for our endpoint. After doing this, we can update the entry for our endpoint in `endpoints.yml` to now state `openapi`.