//Item contatto della lista dei prenotati al mio regalo

import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useBookingStore from "../../../../stores/bookingStore";
import FadeInComponent from "../../../../Components/Atoms/Transitions/FadeInComponent";
import { AsyncImage } from "loadable-image";
import Moment from "react-moment";
import ProfilePic from "../../../../Components/Atoms/ProfilePic/ProfilePic";
import { Booking } from "../../../../core/types/Exchange";
import { User } from "../../../../core/types/User";
import UserAPI from "../../../../core/API/UserAPI";
import DisplayReview from "../../../../Components/Atoms/InLineDisplayers/DisplayReview";
import { Review } from "../../../../core/types/Review";

function ContactItem({ bookingId }: { bookingId: string }) {
    const navigate = useNavigate();
    const [reviews, setReviews] = useState<Review[]>();
    const [loaded, setLoaded] = useState(false);

    const selectedBookingId = useBookingStore(state =>
        state.selectedBooking ? state.selectedBooking._id : null,
    );
    const setSelectedBooking = useBookingStore(state => state.setSelectedBooking);
    const booking: Booking = useBookingStore(state =>
        state.bookings ? state.bookings[bookingId] : null,
    );

    useEffect(() => {
        const handleLoad = async () => {
            const otherUser = await UserAPI.getReviews(booking.otherUser.uid);
            setReviews(otherUser);
            setLoaded(true);
        };

        handleLoad();
    }, [booking.otherUser.uid]);

    return (
        <FadeInComponent>
            {loaded && (
                <div
                    className={`flex flex-col rounded-sm border-2 border-transparent border-b-gray px-4 py-2 ${
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
                    <div className="flex items-center gap-2">
                        {booking.isCanceled ? (
                            <div className="flex items-center gap-2 text-xs">
                                <span className="h-2 w-2 rounded-full bg-gray-dark"></span>
                                <p>Annullato</p>
                            </div>
                        ) : (
                            <Moment locale="it" fromNow className="text-xs text-gray-dark">
                                {booking.createdAt}
                            </Moment>
                        )}
                    </div>
                    <div className="flex items-center gap-3">
                        <ProfilePic userId={booking.otherUser.uid} />
                        <div className="flex flex-col">
                            <p className="semibold">{booking.otherUser.completeName}</p>
                            <DisplayReview reviews={reviews} />
                        </div>
                    </div>
                </div>
            )}
        </FadeInComponent>
    );
}

export default ContactItem;
