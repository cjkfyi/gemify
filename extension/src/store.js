import { createStore } from 'zustand/vanilla';

const store = createStore(() => ({
    // Conversation State 
    conversations: [],
    activeConvoId: null,
    setActiveConvoId: (convoId) => setState({ activeConvoId: convoId }),
    addConversation: (conversation) => setState(state => ({
        conversations: [...state.conversations, conversation]
    })),
    currentResponse: [], // Initialize as an array
    appendResponseChunk: (chunk) => setState(state => ({ currentResponse: [...state.currentResponse, chunk] })),
    setCurrentResponse: (newResponse) => setState({ currentResponse: newResponse }),
}));

const { getState, setState, subscribe, getInitialState } = store;

store.subscribe((newState, previousState) => {
    console.log('Store state changed!', newState);
    // Possibly only return the differences in state? 
    // using previousState, if decided to adjust this
});

export default store;

