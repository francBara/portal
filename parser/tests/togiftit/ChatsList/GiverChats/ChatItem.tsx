import React from "react";
import FadeInComponent from "../../../../Components/Atoms/Transitions/FadeInComponent";
import { AsyncImage } from "loadable-image";
import Moment from "react-moment";
import Location from "../../../../assets/Icons/Location";
import ProfilePic from "../../../../Components/Atoms/ProfilePic/ProfilePic";
import { useNavigate } from "react-router-dom";
import { Gift } from "../../../../core/types/Gift";
import useBookingStore from "../../../../stores/bookingStore";
import { Booking } from "../../../../core/types/Exchange";

interface ChatItemProps {
    gift: Gift;
}

const ChatItem: React.FC<ChatItemProps> = ({ gift }) => {
    const navigate = useNavigate();
    const bookings = useBookingStore(state => state.bookings);
    const keys = Object.keys(bookings);
    const booking = keys.find(key => bookings[key].gift.id === gift._id);
    let giftBookings = [];
    for (let i = 0; i < keys.length; i++) {
        if (bookings[keys[i]].gift.id === gift._id) {
            giftBookings = [...giftBookings, bookings[keys[i]]];
        }
    }
    // Trova il booking sia per regali già consegnati che per regali assegnati
    const bookingRegalato: Booking = giftBookings.find(
        booking =>
            // Controllo per regali già consegnati (confermati da entrambi)
            (booking.isConfirmedByOwner === true && booking.isConfirmedByReceiver === true) ||
            // Controllo per regali assegnati ma non ancora confermati dal ricevente
            booking.isAccepted === true,
    );

    console.log("giftBookings", bookingRegalato);

    return (
        <FadeInComponent>
            <div
                className={`flex items-stretch justify-start rounded-sm border border-b-gray bg-white px-4 py-2`}
                onClick={() => {
                    navigate("/dashboard/" + gift._id);
                }}
            >
                <div className="h-16 w-16 flex-shrink-0 overflow-hidden rounded-lg">
                    {gift?.images?.length ? (
                        <AsyncImage
                            src={gift?.images[0]}
                            className="h-full w-full object-cover"
                            loading="lazy"
                        />
                    ) : (
                        <div className="h-full w-full bg-gray-200" />
                    )}
                </div>
                <div className="w-4" />{" "}
                <div className="flex w-full flex-col justify-between pb-2">
                    <div className="flex w-full justify-between text-sm text-gray-dark">
                        <div className="">
                            <span className="">
                                {" "}
                                <Moment locale="it" fromNow date={gift.createdAt} />
                            </span>
                        </div>
                        <div className="flex items-center">
                            <Location w={20} />
                            <p>{gift?.location?.city}</p>
                        </div>
                    </div>
                    <div className="text-xl font-bold">{gift?.name}</div>
                    <div className="flex items-center">
                        {gift?.isDelivered ? (
                            <div className="flex items-center gap-2">
                                <p>Regalato a</p>
                            </div>
                        ) : gift?.isAssigned ? (
                            <div className="flex items-center gap-2">
                                <p>Assegnato a</p>
                            </div>
                        ) : (
                            <p className="">
                                Prenotato da{" "}
                                <span className="text-verde">{gift.bookingsLength}</span> / 20
                            </p>
                        )}
                        {(gift?.isDelivered || gift?.isAssigned) && (
                            <div className="flex items-center gap-2 pl-2">
                                <ProfilePic
                                    w={6}
                                    key={bookingRegalato.otherUser.uid}
                                    hiddenLvl={true}
                                    image={bookingRegalato.otherUser.image}
                                />

                                <p className="font-semibold underline">
                                    {bookingRegalato.otherUser.completeName.split(" ")[0]}
                                </p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </FadeInComponent>
    );
};

export default ChatItem;
