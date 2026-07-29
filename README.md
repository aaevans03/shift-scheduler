# Shift Scheduler

## Overview

A shift scheduler built in Golang + HTMX.

It is a very simple application. You start as a student user, where you can create and submit a weekly schedule. Then, you can easily switch to an admin user, where you can leave a note and approve or reject the submitted schedule.

All buttons use HTMX interactions to swap in a response from the server. Additionally, the Go backend validates that schedules meet all constraints.

## Running

To run, install Go 1.26.5 and run the following in the project root:

```bash
go run ./cmd/.
```

It will start a local server, in which you will be able to open in your browser at `localhost` + the port provided.
