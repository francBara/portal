import React from "react";

function ChatFilter({ text, chatsFilter, setChatsFilter, index }) {
    const isSelected = chatsFilter === index;

    return (
        <div
            className={`bg-grigio ${isSelected ? "bg-opacity-100" : "bg-opacity-20"} text-${isSelected ? "white" : "black"} cursor-pointer select-none rounded-md px-2 py-2 text-[0.80rem] md:text-sm`}
            onClick={() => {
                setChatsFilter(index);
                if (chatsFilter === index) {
                    setChatsFilter(0);
                }
            }}
        >
            {text}
        </div>
    );
}

export default ChatFilter;
