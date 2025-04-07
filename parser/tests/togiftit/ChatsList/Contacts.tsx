import ContactTile from "./ContactTile";
import ChatItem from "./GiverChats/ChatItem";
import React from "react";
import useBookingStore from "../../../stores/bookingStore";
import { auth } from "../../../firebase";
import Caricamento from "../../../Components/Atoms/Caricamento/Caricamento";
import ChatVuota from "../../../assets/PNG/Chat/ChatVuota.png";
import { Booking } from "../../../core/types/Exchange";
import CardListRegalo from "../../../Components/Atoms/Card/CardListRegalo";
import { Gift } from "../../../core/types/Gift";
import useDeviceDetection from "../../../core/useDeviceDetection";
import { useParams } from "react-router-dom";
import ContactItem from "./GiverChats/ContactItem";
import Logger from "../../../core/logging";
import CardContainer from "../../../Components/Molecules/Container/CardContainer/CardContainer";

function NoChatsText({ text }) {
    return (
        <div className="flex h-full flex-col items-center justify-start gap-5 pt-[10vh] text-grigio md:justify-center md:pt-0">
            {text}
        </div>
    );
}

/**
 *
 * @param {Object} props
 * @param {boolean} props.isReceiver
 * @param {import("../../../core/API/GiftAPI").Gift[]} props.items
 * @returns
 */
function Contacts({ filter, isReceiver, items, list = false }) {
    const isMobile = useDeviceDetection() === "Mobile";
    const urlBookingId = location.pathname.split("/").filter(Boolean).pop();
    /*
     * @type {import("../../../core/API/BookingAPI").Booking[]}
     */

    const bookings: Booking[] = useBookingStore(state => state.bookings);

    return (
        <div className="h-full">
            {(() => {
                if (bookings !== null) {
                    const filteredBookings = Object.values(bookings).filter(booking =>
                        isReceiver
                            ? booking.receiver.uid === auth.currentUser.uid
                            : booking.owner.uid === auth.currentUser.uid,
                    );
                    console.log("Filtered bookings: ", filteredBookings);

                    if (bookings.length === 0) {
                        return <NoChatsText text={""} />;
                    }

                    let filteredChats = filteredBookings
                        .sort((a, b) => {
                            // First, put canceled bookings at the bottom
                            if (a.isCanceled && !b.isCanceled) {
                                return 1;
                            }
                            if (!a.isCanceled && b.isCanceled) {
                                return -1;
                            }

                            // Then sort by last message timestamp
                            if (a.lastMessage === null) {
                                return -1;
                            } else if (b.lastMessage === null) {
                                return 1;
                            }
                            return (
                                b.lastMessage.createdAt.getTime() -
                                a.lastMessage.createdAt.getTime()
                            );
                        })
                        // Chats are filtered if a filter is selected
                        .filter(filter);

                    if (items !== null && items !== undefined) {
                        if (items?.length === 0) {
                            return (
                                <div className="">
                                    {
                                        <img
                                            src={ChatVuota}
                                            className="mx-auto mt-10 w-[65vw] md:hidden"
                                            alt=""
                                        />
                                    }
                                    <NoChatsText
                                        text={"Qui vedrai le chat dei tuoi regali pubblicati"}
                                    />
                                </div>
                            );
                        }

                        const groupedChats = {};

                        for (let item of items) {
                            groupedChats[item._id] = [];
                        }

                        for (let filteredChat of filteredChats) {
                            if (filteredChat.gift.id in groupedChats) {
                                groupedChats[filteredChat.gift.id].push(filteredChat);
                            } else {
                                groupedChats[filteredChat.gift.id] = [filteredChat];
                            }
                        }
                        const currentGift = bookings[urlBookingId]?.gift.id;
                        filteredChats = filteredChats.filter(chat => chat.gift.id === currentGift);

                        if (isMobile) {
                            return items.map((item: Gift, index: number) => {
                                return (
                                    <div className="" key={index}>
                                        <ChatItem gift={item} />
                                    </div>
                                );
                            });
                        }
                        if (list) {
                            if (filteredChats.length > 0) {
                                return (
                                    <div className="h-full overflow-y-auto bg-white">
                                        {filteredChats.map((chat: Booking, index: number) => {
                                            return (
                                                <ContactItem bookingId={chat._id} key={chat._id} />
                                            );
                                        })}
                                    </div>
                                );
                            }
                        }
                        return (
                            <div className="flex h-full justify-center bg-gray-ultralight pt-10">
                                <CardContainer cols={4} grid>
                                    {items.map((item: Gift) => {
                                        return <CardListRegalo gift={item} />;
                                    })}
                                </CardContainer>
                            </div>
                        );
                    }

                    if (filteredChats.length > 0) {
                        return (
                            <div className="overflow-y-auto rounded-b-lg bg-white">
                                {filteredChats.map(chat => {
                                    return <ContactTile bookingId={chat._id} key={chat._id} />;
                                })}
                            </div>
                        );
                    }
                    return (
                        <div className="">
                            {
                                <img
                                    src={ChatVuota}
                                    className="mx-auto mt-10 w-[65vw] md:hidden"
                                    alt=""
                                />
                            }
                            <NoChatsText text={"Qui vedrai le chat delle tue prenotazioni"} />
                        </div>
                    );
                }
                return <Caricamento />;
            })()}
        </div>
    );
}

export default Contacts;
