import ProfilePic from "../../../Components/Atoms/ProfilePic/ProfilePic";
import FadeInComponent from "../../../Components/Atoms/Transitions/FadeInComponent";
import { auth } from "../../../firebase";
import React, { useEffect, useState } from "react";
import useBookingStore from "../../../stores/bookingStore";
import { useNavigate } from "react-router-dom";
import GiftAPI from "../../../core/API/GiftAPI";
import { AsyncImage } from "loadable-image";
import Moment from "react-moment";
import Location from "../../../assets/Icons/Location";
import { Booking } from "../../../core/types/Exchange";

/**
 *
 * @param {Object} props
 * @param {string} props.bookingId
 * @returns
 */

function ContactTile({ bookingId }: { bookingId: string }) {
    const navigate = useNavigate();
    const [gift, setGift] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    const selectedBookingId = useBookingStore(state =>
        state.selectedBooking ? state.selectedBooking._id : null,
    );
    const setSelectedBooking = useBookingStore(state => state.setSelectedBooking);
    const booking: Booking = useBookingStore(state =>
        state.bookings ? state.bookings[bookingId] : null,
    );

    useEffect(() => {
        const loadGiftImage = async () => {
            if (booking?.gift?.id) {
                try {
                    setIsLoading(true);
                    const res = await GiftAPI.getById(booking.gift.id);

                    setGift(res);
                } catch (error) {
                    console.error("Errore nel caricamento dell'immagine:", error);
                } finally {
                    setIsLoading(false);
                }
            }
        };

        loadGiftImage();
    }, [booking?.gift?.id]);

    if (!booking) return null;

    return (
        <FadeInComponent>
            <div
                className={`flex items-stretch justify-start rounded-sm border-2 border-transparent border-b-gray px-4 py-4 ${
                    selectedBookingId === booking._id
                        ? "bg-verde/30"
                        : true //booking.isRead
                          ? "bg-white"
                          : "bg-gray-light"
                } cursor-pointer transition duration-300`}
                onClick={() => {
                    navigate("/chat/" + booking._id);
                    setSelectedBooking(booking);
                }}
            >
                <div className="h-16 w-16 flex-shrink-0 overflow-hidden rounded-lg">
                    {gift ? (
                        <AsyncImage
                            src={gift?.images[0]}
                            className="h-full w-full object-cover"
                            loading="lazy"
                        />
                    ) : (
                        <div className="h-full w-full bg-gray-200" />
                    )}
                </div>
                <div className="w-4" />
                <div className="flex w-full flex-col justify-between">
                    <div className="flex w-full justify-between text-sm text-gray-dark">
                        <div className="">
                            {booking.isCanceled ? (
                                <div className="flex items-center gap-2">
                                    <span className="h-2 w-2 rounded-full bg-gray-dark"></span>
                                    <p>Annullato</p>
                                </div>
                            ) : gift?.isDelivered ? (
                                <div className="flex items-center gap-2">
                                    <span className="h-2 w-2 rounded-full bg-verde"></span>
                                    <p>Regalato</p>
                                </div>
                            ) : gift?.isAssigned ? (
                                <div className="flex items-center gap-2">
                                    <span className="h-2 w-2 rounded-full bg-warning"></span>
                                    <p>Assegnato</p>
                                </div>
                            ) : (
                                <span className="">
                                    {" "}
                                    <Moment locale="it" fromNow date={booking.createdAt} />
                                </span>
                            )}
                        </div>
                        <div className="flex items-center">
                            <Location w={20} />
                            <p>{gift?.location?.city}</p>
                        </div>
                    </div>
                    <div className="text-xl font-bold">{booking.gift?.name}</div>
                    <div className="flex items-center gap-1 text-base">
                        <p className="mr-2">di</p>
                        <ProfilePic
                            w={6}
                            key={booking.otherUser.uid}
                            hiddenLvl={true}
                            image={booking.otherUser.image}
                        />

                        <p className="font-semibold underline">{booking.otherUser?.completeName}</p>
                    </div>
                </div>
            </div>
        </FadeInComponent>
    );
}

export default ContactTile;
