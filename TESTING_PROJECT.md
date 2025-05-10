I'll create a concise specification for a markdown-based ticket tracking system for Claude Code to implement. This will help maintain context between sessions and provide clear prioritization of tasks.

Here's what I'd recommend:

```
# Task Tracking System Specification for Claude Code

## Overview
Implement a self-maintaining markdown task tracking system to maintain context between sessions. This system will focus on test coverage improvement tasks.

## Core Requirements
1. Single markdown file to track all tasks
2. Status tracking: "TODO", "DOING", "DONE"
3. Task workflow:
   - Start with "DOING" tasks first
   - When picking up a "TODO" task, update to "DOING"
   - When completing a "DOING" task, update to "DONE"
   - "Done" is defined by successful test completion

## File Structure
```markdown
# Task Tracking System

## Active Tasks
<!-- Tasks currently in progress have DOING status -->

- [DOING] Implement unit tests for authentication module (#1)
  - [ ] Test user login flow
  - [ ] Test token validation
  - Progress: 1/3 subtasks complete

## Backlog
<!-- Tasks not yet started have TODO status -->

- [TODO] Add integration tests for payment processing (#2)
  - [ ] Test successful payment flow
  - [ ] Test declined payment scenarios
  - [ ] Test refund process

## Completed
<!-- Tasks that are finished have DONE status -->

- [DONE] Set up testing framework (#0)
  - [x] Install test dependencies
  - [x] Configure CI pipeline
  - [x] Create test directory structure
  - Completed: 2023-05-08, All tests passing
```

## Implementation Logic
1. On startup:
   - Parse the markdown file to identify all tasks
   - Process tasks in this order: DOING → TODO
   - Never start a new TODO task if a DOING task exists

2. For task management:
   - When starting a TODO task:
     - Move to Active Tasks section
     - Change status to [DOING]
   - When completing a DOING task:
     - Move to Completed section
     - Change status to [DONE]
     - Add completion date and test status
   - When creating a new task:
     - Add to Backlog with [TODO] status
     - Assign unique ID number

3. Task properties:
   - Unique ID (#number)
   - Status (TODO/DOING/DONE)
   - Title (brief description)
   - Subtasks (checkable items)
   - Metadata (dates, test status)

Claude Code should maintain this file automatically, updating task status and location as work progresses. This will provide a persistent record of project status between sessions.
```

This structure will give Claude Code a clear framework for tracking test coverage improvement tasks while maintaining state between sessions. It includes all the key functionality you requested: status tracking, proper task workflow, and a way to track task completion through test results.
