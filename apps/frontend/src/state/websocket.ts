// import { create } from "zustand";
// import WebSocketManager, {
//   Subscription,
//   WebSocketMessage,
// } from "../services/websocket";

// interface WebSocketState {
//   isConnected: boolean;
//   messages: WebSocketMessage[];
//   wsUtility: WebSocketManager | null;
//   connect: (url: string) => void;
//   subscribe: (subscription: Subscription) => void;
//   unsubscribe: (subscription: Subscription) => void;
//   close: () => void;
// }

// const useWebSocketStore = create<WebSocketState>((set, get) => ({
//   isConnected: false,
//   messages: [],
//   wsUtility: null,

//   connect: (url: string) => {
//     const wsUtility = new WebSocketManager(url);

//     wsUtility.onMessage((message) => {
//       set((state) => ({ messages: [...state.messages, message] }));
//     });

//     wsUtility.connect();
//     set({ wsUtility, isConnected: true });
//   },

//   subscribe: (subscription: Subscription) => {
//     const { wsUtility } = get();
//     if (wsUtility) {
//       wsUtility.subscribe(subscription);
//     }
//   },

//   unsubscribe: (subscription: Subscription) => {
//     const { wsUtility } = get();
//     if (wsUtility) {
//       wsUtility.unsubscribe(subscription);
//     }
//   },

//   close: () => {
//     const { wsUtility } = get();
//     if (wsUtility) {
//       wsUtility.close();
//       set({ wsUtility: null, isConnected: false });
//     }
//   },
// }));

// export default useWebSocketStore;
