import { createStore } from 'zustand/vanilla';

const store = createStore(() => ({
    proj: null,
    setProj: (obj) => setState({ proj: obj }),
    chat: null,
    setChat: (obj) => setState({ chat: obj }),
}));

const { 
    getState, 
    setState, 
    subscribe, 
    getInitialState 
} = store;

store.subscribe((newState, previousState) => {
    console.log('Store state has changed!', newState);
});

export default store;

