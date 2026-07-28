/**
 * SCHEDULE CREATION
 * 1. On page open, grab schedule from server with HTMX
 *      -> load it into the scheduler frame
 * 2. Press edit button -> a flow starts.
 * 3. Whatever schedule edits you make is stored in immediate memory
 * 4. Submit -> info is sent to the server
 */

function initializeScheduleEditor(root = document) {

    const inputContainer = root.querySelector('#form-selected-blocks');
    if (!inputContainer) return;
    
    function updateSelectedInputs() {
        inputContainer.innerHTML = '';
        
        root.querySelectorAll('.block.active').forEach((block) => {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = block.dataset.day;
            input.value = block.dataset.time;
            inputContainer.appendChild(input);
        })
    }
    
    function toggleBlock(block) {
        if (toggleOn === true) {
            block.classList.add('active');
        } else {
            block.classList.remove('active');
        }
        updateSelectedInputs();
    }
    
    
    let isPainting = false;
    let toggleOn = null;
    const allBlocks = root.querySelectorAll('.edit');
    
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
    
    root.addEventListener('pointerup', () => {
        isPainting = false;
        toggleOn = null;
    }, { once: false });
    
    updateSelectedInputs();
}

document.addEventListener('DOMContentLoaded', () => {
    initializeScheduleEditor();
});

document.body.addEventListener('htmx:afterSwap', (event) => {
    initializeScheduleEditor(event.target);
});
