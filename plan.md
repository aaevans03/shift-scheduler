# Shift Scheduler

What is my attack plan here?

## Steps

- [X] Start with static HTML
    - Figure out all of the basic layout. How are things layered?
- [X] Have a Go server serve the static files
- [X] Start splitting things up into HTMX reusable templates
    - [X] Layered components
- [ ] HTMX Core interactions, constraints, edge cases
    - [X] Toggle each square on and off
        - Is it overkill to send every toggle to the server to do a time calculation? I think so.
    - [X] Replace table content dynamically
        - On page open: Read from DB, send a schedule object back
    - [X] Submit the form without a page reload
    - [X] Time validation
    - [X] Rendering server-side validation errors
        - Like error pop-ups
    - [ ] View versus edit mode
    - [ ] Status updates (approved/rejected banners, etc.)
- [ ] Admin
