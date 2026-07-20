# Shift Scheduler

What is my attack plan here?

## Steps

- [X] Start with static HTML
    - Figure out all of the basic layout. How are things layered?
- [X] Have a Go server serve the static files
- [ ] Start splitting things up into HTMX reusable templates
    - [ ] Layered components
- [ ] HTMX Core interactions, constraints, edge cases
    - [ ] Toggle each square on and off
        - Is it overkill to send every toggle to the server to do a time calculation? I think so.
    - [ ] Replace table content dynamically
        - On page open: Read from DB, send a schedule object back
    - [ ] Submit the form without a page reload
    - [ ] Rendering server-side validation errors
        - Like error pop-ups
    - [ ] Status updates (approved/rejected banners, etc.)
