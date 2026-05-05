// Filename: static/js/core/event-emitter.js

function createEmitter() {
    // Map to store event listeners: key = event name, value = array of callback functions
    const listeners = new Map();

    // Register an event listener
    function on(event, callback) {
        if (!listeners.has(event)) {
            listeners.set(event, [])
        }
        listeners.get(event).push(callback)
    }

    // Trigger an event and notify all listeners
    function emit(event, data) {
        if (!listeners.has(event)) {
            return;
        }
        listeners.get(event).forEach(cb => cb(data));
    }

    function off(event, callback) {
        if (!listeners.has(event)) {
            return;
        }
        listeners.set(
            event,
            listeners.get(event).filter(cb => cb !== callback)
        );
    }

    return { on, emit, off }
}

// Create a single global emitter instance
export const emitter = createEmitter();