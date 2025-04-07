import React, { useEffect, useState } from "react";
import ReceiverChats from "./ReceiverChats";
import GiverChats from "./GiverChats/GiverChats";
import useBookingStore from "../../../stores/bookingStore";
import BookingAPI, { parseBooking } from "../../../core/API/BookingAPI";
import { useNavigate, useParams } from "react-router-dom";
import { auth } from "../../../firebase";
import useNavbarStore from "../../../stores/navbarStore";
import { motion } from "framer-motion";
import { getAbsoluteImage } from "../../../core/API/CallBackend";
import Header from "../../../Components/Molecules/Sidebar/Header";
import ChipTabs from "../../../Components/Atoms/Tabs/ChipTabs";

function ChatsList({ showTitle = false, list = false }) {
    const receiverSelected = useBookingStore(state => state.receiverSelected);
    const setReceiverSelected = useBookingStore(state => state.setReceiverSelected);
    const selectedBooking = useBookingStore(state => state.selectedBooking);
    const setSelectedBooking = useBookingStore(state => state.setSelectedBooking);

    const setNavbarVisible = useNavbarStore(state => state.setNavbarVisible);

    const navigate = useNavigate();

    const [selected, setSelected] = useState("Prenotazioni");

    const handleSelect = tab => {
        if (tab === "Prenotazioni") {
            if (!receiverSelected) {
                setReceiverSelected(true);
                setSelectedBooking({});
            }
            navigate("/chat");
        } else {
            setReceiverSelected(false);
            setSelectedBooking({});
            navigate("/chat");
        }
        setSelected(tab);
    };

    useEffect(() => {
        if (receiverSelected) {
            setSelected("Prenotazioni");
        } else {
            setSelected("Donazioni");
        }
    }, [receiverSelected]);

    return (
        <div className="flex h-full flex-col bg-gray-ultralight">
            {showTitle && (
                <Header
                    title="Trattative"
                    isSubHeader
                    hideBackButton
                    backgroundType="transparent"
                />
            )}
            {!list && (
                <div className={`${receiverSelected || list ? "w-full" : "md:w-1/2 lg:w-1/3"}`}>
                    <ChipTabs
                        selected={selected}
                        setSelected={e => handleSelect(e)}
                        tabs={["Prenotazioni", "Donazioni"]}
                    />
                </div>
            )}
            <div className="flex h-full flex-grow flex-col overflow-auto">
                {receiverSelected ? <ReceiverChats /> : <GiverChats list={list} />}
            </div>
        </div>
    );
}

export default ChatsList;

/*
<button className={`mr-4 ml-8 font-bold ${receiverSelected ? "text-black" : "text-gray-400"}`} onClick={() => {
                    if (!receiverSelected) {
                        navigate('/chat/');
                        setSelectedBooking({});
                    }
                    setReceiverSelected(true);
                }}>
                    Adozioni
                </button>
                <button className={`font-bold ${receiverSelected ? "text-gray-400" : "text-black"}`} onClick={() => {
                    if (receiverSelected) {
                        navigate('/chat/');
                        setSelectedBooking({});
                    }
                    setReceiverSelected(false);
                }}>
                    Donazioni
                </button>
*/
