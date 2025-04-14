import { useEffect, useState } from "react";
import ChatMessages from "./ChatMessages/ChatMessages";
import ChatsList from "./ChatsList/ChatsList";
import useBookingStore from "../../stores/bookingStore";
import { useLocation } from "react-router-dom";
import BookingAPI, { parseBooking } from "../../core/API/BookingAPI";
import { auth } from "../../firebase";
import { getAbsoluteImage } from "../../core/API/CallBackend";
import ChatVuota from "../../assets/PNG/Chat/ChatVuota.png";
import Logger from "../../core/logging";
import BannerRegalo from "./Components/BannerRegalo";
import { toast } from "sonner";
import Indietro from "../../Components/Atoms/Bottoni/Indietro";

function Chat() {
    const bookings = useBookingStore(state => state.bookings);
    const setBookings = useBookingStore(state => state.setBookings);

    const updateBookings = useBookingStore(state => state.updateBookings);
    const appendMessage = useBookingStore(state => state.appendMessage);

    const setSelectedBooking = useBookingStore(state => state.setSelectedBooking);
    const selectedBooking = useBookingStore(state => state.selectedBooking);
    const receiverSelected = useBookingStore(state => state.receiverSelected);
    const setReceiverSelected = useBookingStore(state => state.setReceiverSelected);

    //TODO: Ask Andrea if it is correct
    const [isMobile, setIsMobile] = useState(window.innerWidth < 990);

    const location = useLocation();

    //@portal

    //@portal
    const maxChats = 24;

    //@portal

    let asd = 2;

    // @portal
    let chatName = "My Chat";

    const handleUrl = (bookings: Record<string, any>) => {
        const urlBookingId = location.pathname.split("/").filter(Boolean).pop();

        if (urlBookingId) {
            if (urlBookingId in bookings) {
                if (bookings[urlBookingId].owner.uid === auth.currentUser?.uid) {
                    setReceiverSelected(false);
                }

                setSelectedBooking(bookings[urlBookingId]);
            } else {
                setSelectedBooking({});
            }
        }
    };

    const loadBookings = async () => {
        try {
            const bookings = await BookingAPI.getBookings();

            const mappedBookings: Record<string, any> = {};

            for (let booking of bookings.flat()) {
                mappedBookings[booking._id] = booking;
            }

            handleUrl(mappedBookings);

            setBookings(mappedBookings);

            const closeConnection = BookingAPI.onBookingUpdate(
                //On booking status change (canceled, confirmed ecc.)
                (booking: any) => {
                    booking = parseBooking(booking);
                    updateBookings(booking);
                },
                //On new chat message, message contains bookingId
                message => {
                    if (message.image) {
                        message.image = getAbsoluteImage(message.image);
                    }
                    appendMessage(message);
                },
            );

            return () => {
                closeConnection();
            };
        } catch (e) {
            Logger.error("Error loading bookings", e);
            toast.error("Errore durante il caricamento delle chat");
        }
    };

    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const donazioni = params.get("donazioni");
        if (donazioni !== null) {
            setReceiverSelected(false);
        }

        const handleResize = () => {
            setIsMobile(window.innerWidth < 990);
        };

        window.addEventListener("resize", handleResize);

        return () => {
            window.removeEventListener("resize", handleResize);
        };
    }, []);

    useEffect(() => {
        if (bookings === null) {
            loadBookings();
        } else {
            handleUrl(bookings);
        }
    }, [location]);

    return (
        <div className="min-h-[90vh] bg-gray-ultralight">
            <div className="mx-auto md:w-[70%] md:py-5">
                <div className="absolute left-3 flex items-center">
                    {!receiverSelected && !isMobile && selectedBooking._id && <Indietro />}
                </div>
                <h2 className="mb-3 hidden text-4xl font-bold md:block">Trattative</h2>
                {!receiverSelected &&
                    selectedBooking !== null &&
                    selectedBooking._id &&
                    !isMobile && (
                        <div className="mb-5">
                            <BannerRegalo booking={selectedBooking} />
                        </div>
                    )}
                {isMobile ? (
                    <div className="h-[100vh] rounded-lg bg-white">
                        {selectedBooking._id ? (
                            <ChatMessages key={selectedBooking._id} />
                        ) : (
                            <div className="h-full rounded-lg bg-white">
                                <ChatsList />
                            </div>
                        )}
                    </div>
                ) : !receiverSelected && !isMobile ? (
                    <div className="h-[80vh] rounded-lg">
                        {selectedBooking._id ? (
                            <div className="flex h-[70vh] justify-around">
                                <div className="h-full rounded-lg bg-white shadow md:w-1/2 lg:w-1/3">
                                    <ChatsList list />
                                </div>
                                <div className="w-4"></div>
                                <div className="h-full rounded-lg bg-white shadow md:w-1/2 lg:w-2/3">
                                    {selectedBooking._id ? (
                                        <ChatMessages key={selectedBooking._id} />
                                    ) : (
                                        <div className="flex h-full flex-col items-center justify-center gap-5 text-center text-grigio">
                                            <img
                                                src={ChatVuota}
                                                className="hidden w-[35vw] md:block"
                                                alt=""
                                            />
                                            <p>Nessuna chat selezionata</p>
                                        </div>
                                    )}
                                </div>
                            </div>
                        ) : (
                            <div className="h-full rounded-lg bg-white">
                                <ChatsList />
                            </div>
                        )}
                    </div>
                ) : (
                    <div className="flex h-[80vh] justify-around">
                        <div className="h-full rounded-lg bg-white shadow md:w-1/2 lg:w-1/3">
                            <ChatsList />
                        </div>
                        <div className="w-4"></div>
                        <div className="h-full rounded-lg bg-white shadow md:w-1/2 lg:w-2/3">
                            {selectedBooking._id ? (
                                <ChatMessages key={selectedBooking._id} />
                            ) : (
                                <div className="flex h-full flex-col items-center justify-center gap-5 text-center text-grigio">
                                    <img
                                        src={ChatVuota}
                                        className="hidden w-[35vw] md:block"
                                        alt=""
                                    />
                                    <p>Nessuna chat selezionata</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}

export default Chat;
