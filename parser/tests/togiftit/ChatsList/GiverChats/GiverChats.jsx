import { useEffect, useState } from "react";
import ChatFilter from "../ChatFilter";
import Contacts from "../Contacts";
import { auth } from "../../../../firebase";
import React from "react";
import GiftAPI from "../../../../core/API/GiftAPI";
import Caricamento from "../../../../Components/Atoms/Caricamento/Caricamento";

function GiverChats({ list = false }) {
    const [chatsFilter, setChatsFilter] = useState(0);
    /**
     * @type {[import("../../../../core/API/GiftAPI").Gift[], function]}
     */
    const [items, setItems] = useState(null);

    /** @param {import("../../../../core/API/BookingAPI").Booking} chat **/
    const nullFilter = chat => {
        return [!chat.isCanceled, chat.isCanceled];
    };

    /** @param {import("../../../../core/API/BookingAPI").Booking} chat **/
    const infoFilter = chat => {
        return !chat.isRequested && !chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../../core/API/BookingAPI").Booking} chat **/
    const bookedFilter = chat => {
        return chat.isRequested && !chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../../core/API/BookingAPI").Booking} chat **/
    const acceptedFilter = chat => {
        return chat.isRequested && chat.isAccepted && !chat.isCanceled;
    };

    /** @param {import("../../../../core/API/BookingAPI").Booking} chat **/
    const canceledFilter = chat => {
        return chat.isCanceled;
    };

    const filters = [nullFilter, infoFilter, bookedFilter, acceptedFilter, canceledFilter];

    const getItems = async () => {
        const gifts = await GiftAPI.getByOwner(auth.currentUser.uid);
        setItems(gifts);
    };

    useEffect(() => {
        if (items === null) {
            getItems();
        }
    }, []);

    return (
        <div className="flex h-full flex-col">
            <div className="flex-grow overflow-y-auto">
                {items !== null ? (
                    <Contacts
                        filter={filters[chatsFilter]}
                        isReceiver={false}
                        list={list}
                        items={items}
                        key={"giver"}
                    />
                ) : (
                    <Caricamento />
                )}
            </div>
        </div>
    );
}

export default GiverChats;
