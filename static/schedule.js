/**
 * SCHEDULE CREATION
 * 1. On page open, grab schedule from server with HTMX
 *      -> load it into the scheduler frame
 * 2. Press edit button -> a flow starts.
 * 3. Whatever schedule edits you make is stored in immediate memory
 * 4. Submit -> info is sent to the server
 */

let isPainting = false;
let toggleOn = null;

const allBlocks = document.querySelectorAll('.block');

function toggleBlock(block) {
    if (toggleOn === true) {
        block.classList.add('active');
    } else {
        block.classList.remove('active');
    }
}

allBlocks.forEach((block) => {
    block.addEventListener('pointerdown', (event) => {

        isPainting = true;

        if (event.currentTarget.classList.contains('active')) {
            toggleOn = false;
        } else toggleOn = true;
        
        toggleBlock(event.currentTarget);
    });

    block.addEventListener('pointerenter', (event) => {
        if (!isPainting) return;
        toggleBlock(event.currentTarget);
    });
});

document.addEventListener('pointerup', () => {
    isPainting = false;
    toggleOn = null;
})
