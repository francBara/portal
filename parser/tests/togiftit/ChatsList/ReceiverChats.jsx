import { useState } from "react";
import ChatFilter from "./ChatFilter";
import Contacts from "./Contacts";
import React from "react";

function ReceiverChats() {
    const [chatsFilter, setChatsFilter] = useState(0);

    /** @param {import("../../../core/API/BookingAPI").Booking} chat **/
    const nullFilter = chat => {
        return [!chat.isCanceled, chat.isCanceled];
    };

    /** @param {import("../../../core/API/BookingAPI").Booking} chat **/
    const infoFilter = chat => {
        return !chat.isRequested && !chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../core/API/BookingAPI").Booking} chat **/
    const bookedFilter = chat => {
        return chat.isRequested && !chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../core/API/BookingAPI").Booking} chat **/
    const acceptedFilter = chat => {
        return chat.isRequested && chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../core/API/BookingAPI").Booking} chat **/
    const canceledFilter = chat => {
        return chat.isCanceled;
    };

    const filters = [nullFilter, infoFilter, bookedFilter, acceptedFilter, canceledFilter];

    return (
        <div className="flex h-full flex-col">
            <div className="mb-4 ml-8 hidden flex-row">
                <ChatFilter
                    text={"Contattati"}
                    chatsFilter={chatsFilter}
                    setChatsFilter={setChatsFilter}
                    index={1}
                />
                <div className="w-3" />
                <ChatFilter
                    text={"Prenotati"}
                    chatsFilter={chatsFilter}
                    setChatsFilter={setChatsFilter}
                    index={2}
                />
                <div className="w-3" />
                <ChatFilter
                    text={"Assegnati"}
                    chatsFilter={chatsFilter}
                    setChatsFilter={setChatsFilter}
                    index={3}
                />
                <div className="w-3" />
                <ChatFilter
                    text={"Annullati"}
                    chatsFilter={chatsFilter}
                    setChatsFilter={setChatsFilter}
                    index={4}
                />
            </div>
            <div className="flex-grow overflow-y-auto">
                <Contacts filter={filters[chatsFilter]} isReceiver={true} key={"receiver"} />
            </div>
        </div>
    );
}

export default ReceiverChats;
