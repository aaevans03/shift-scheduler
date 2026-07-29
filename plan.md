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
    - [X] View versus edit mode
    - [X] Status updates (approved/rejected banners, etc.)
    - [X] Daily and weekly hour totals
- [ ] Admin approval functionality
    - [X] Button to switch to admin
        - Persist a "logged in" user with cookies or session map, with hardcoded user roles.
    - [ ] Approve and reject submissions, with comments
    - [ ] Approved schedule is read-only unless reset
    - [ ] Users can edit & resubmit if previously rejected
