const chatState = {
    wasChatParamSelected: false,
    isReceiverSectionSelected: true,
};

function initChatState() {
    chatState.wasChatParamSelected = false;
    chatState.isReceiverSectionSelected = true;
}

export { chatState, initChatState };
