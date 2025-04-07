import React, { useEffect, useRef, useState, useContext } from "react";
import { auth } from "../../../firebase";
import ProfilePic from "../../../Components/Atoms/ProfilePic/ProfilePic";
import Caricamento from "../../../Components/Atoms/Caricamento/Caricamento";
import RenderMessage from "./RenderMessage";
import ConfirmBookingTile from "./ConfirmBookingTile";
import Iconfrecciagiu from "../../../assets/Icons/Iconfrecciagiu";
import { ReducerContext } from "../../../rootReducer";
import ChatAPI from "../../../core/API/ChatAPI";
import useBookingStore from "../../../stores/bookingStore";
import { useNavigate } from "react-router-dom";
import chatBackground from "../../../assets/Sfondi/chat_texture2.svg";
import chatBackgroundDesktop from "../../../assets/Textures/texture_square_grigio_2.png";
import { toast } from "sonner";
import Camera from "../../../assets/Icons/Camera";
import { getAbsoluteImage } from "../../../core/API/CallBackend";
import { acceptedImageTypes } from "../../../constants/fileExtensions";
import Indietro from "../../../Components/Atoms/Bottoni/Indietro";
import DefaultButton from "../../../Components/Atoms/Bottoni/DefaultButton";
import Send from "../../../assets/Icons/Send";
import BannerRegalo from "../Components/BannerRegalo";
import { Message } from "../../../core/types/Chat";

function ChatMessages() {
    const [text, setText] = useState("");
    const loaded = useRef(false);
    const [imagePreview, setImagePreview] = useState(null);
    const { dispatch } = useContext(ReducerContext);
    const isMobile = /Mobi|Android/i.test(navigator.userAgent);
    const backgroundImage = isMobile ? chatBackground : chatBackgroundDesktop;

    const navigate = useNavigate();

    /**
     * @type {import("../../../core/API/ChatAPI").Message[]}
     */
    const messages = useBookingStore(state => state.displayedMessages);
    const setMessages = useBookingStore(state => state.setDisplayedMessages);
    const appendMessage = useBookingStore(state => state.appendMessage);

    const scrollRef = useRef(null);
    const fixedDivRef = useRef(null);

    const fileInputRef = useRef(null);
    const textInputRef = useRef(null);
    let tmpImageFile = useRef(null);

    /**
     * @type {import("../../../core/API/BookingAPI").Booking}
     */
    const booking = useBookingStore(state => state.selectedBooking);
    const setSelectedBooking = useBookingStore(state => state.setSelectedBooking);
    const updateBooking = useBookingStore(state => state.updateBookings);

    useEffect(() => {
        if (!loaded.current && booking._id) {
            dispatch({ type: "hide_sidebar_phone" });
            load();
        }
    }, [booking._id]);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, []);

    const load = async () => {
        const chat = await ChatAPI.getMessages(booking._id);

        loaded.current = true;

        setMessages(chat.messages.reverse());
        updateBooking({
            _id: booking._id,
            isRead: true,
        });

        if (scrollRef.current) {
            scrollRef.current.scrollTop = 0;
        }
    };

    const onPhotoChoose = e => {
        tmpImageFile.current = e.target.files[0];
        const imageUrl = URL.createObjectURL(tmpImageFile.current);
        setImagePreview(imageUrl);
        const fileInput = e.target;
        fileInput.value = null;
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
        if (textInputRef.current) {
            textInputRef.current.focus();
        }
    };

    function checkText(text) {
        return text.trim().length > 0;
    }

    const send = async e => {
        e.preventDefault();

        const content = text;
        const image = tmpImageFile.current;
        const tmpImagePreview = imagePreview;

        tmpImageFile.current = null;
        setText("");
        setImagePreview(null);

        if ((!checkText(content) && !image) || booking.isBlocked) {
            return;
        }

        if (scrollRef.current) {
            scrollRef.current.scrollTop = 0;
        }

        try {
            let imageUrl = null;

            const randomId = Math.floor(Math.random() * 10000);

            appendMessage({
                content: content,
                sendBy: auth.currentUser.uid,
                createdAt: new Date(),
                image: tmpImagePreview,
                bookingId: booking._id,
                isLoading: true,
                id: randomId,
            });

            if (image) {
                imageUrl = await ChatAPI.uploadMessageImage(image, booking._id);
            }

            const completeImageUrl = imageUrl ? getAbsoluteImage(imageUrl) : null;

            let sent = false;

            for (let i = 0; i < 3; i++) {
                try {
                    await ChatAPI.sendMessage(booking._id, {
                        content: content,
                        image: imageUrl,
                    });
                    sent = true;
                    break;
                } catch (e) {
                    if (!e.status || e.status !== 429) {
                        throw e;
                    }
                }
            }

            if (!sent) {
                throw "Message not sent";
            }

            appendMessage({
                content: content,
                sendBy: auth.currentUser.uid,
                createdAt: new Date(),
                image: completeImageUrl,
                bookingId: booking._id,
                id: randomId,
            });
        } catch (e) {
            toast.error("Si è verificato un errore durante la consegna del messaggio");
        }
    };

    let fired = false;

    window.onkeydown = function (e) {
        if (!fired && e.key === "Enter") {
            fired = true;
            if (checkText(text)) {
                send(e);
            }
        }
    };

    window.onkeyup = function () {
        fired = false;
    };

    useEffect(() => {
        const scrollableDiv = scrollRef.current;

        if (!isMobile || !scrollableDiv || !fixedDivRef.current) return;

        let touchStarted = false;
        let scrolling = false;
        let scrollTimeout;

        const hideFixedDiv = () => {
            if (fixedDivRef.current) {
                fixedDivRef.current.style.transform = "translateY(-150%)";
                fixedDivRef.current.style.transition = "transform 0.3s ease";
            }
        };

        const showFixedDiv = () => {
            if (fixedDivRef.current) {
                fixedDivRef.current.style.transform = "translateY(0)";
                fixedDivRef.current.style.transition = "transform 0.3s ease";
            }
        };

        const handleTouchStart = () => {
            touchStarted = true;
            clearTimeout(scrollTimeout);
        };

        const handleTouchMove = () => {
            if (touchStarted && !scrolling) {
                scrolling = true;
                hideFixedDiv();
            }
        };

        const handleTouchEnd = () => {
            touchStarted = false;
            scrolling = false;
            clearTimeout(scrollTimeout);
            scrollTimeout = setTimeout(() => {
                showFixedDiv();
            }, 500);
        };

        scrollableDiv.addEventListener("touchstart", handleTouchStart);
        scrollableDiv.addEventListener("touchmove", handleTouchMove);
        scrollableDiv.addEventListener("touchend", handleTouchEnd);

        return () => {
            scrollableDiv.removeEventListener("touchstart", handleTouchStart);
            scrollableDiv.removeEventListener("touchmove", handleTouchMove);
            scrollableDiv.removeEventListener("touchend", handleTouchEnd);
            clearTimeout(scrollTimeout);
        };
    }, [isMobile]);

    if (!booking._id) {
        return <div />;
    }

    return (
        <div className={`${isMobile && "pb-16"} flex h-full flex-col`}>
            <div
                className={`${
                    isMobile && "fixed left-0 top-0 w-full"
                } z-10 flex h-12 flex-col items-center border-b-2 border-b-grigino bg-white text-base font-light shadow-sm md:h-24`}
            >
                <div className="relative z-20 flex w-full items-center bg-white p-4 shadow-sm md:mb-2">
                    <span className="md:hidden">
                        <Indietro
                            customWidth={32}
                            onBeforeClick={() => dispatch({ type: "show_sidebar_phone" })}
                        />
                    </span>
                    <div className="flex w-full cursor-pointer items-center justify-center text-right font-light md:flex-row-reverse md:justify-end md:text-left">
                        <div className="mr-3 flex flex-col gap-1 md:ml-3">
                            <p
                                onClick={() => navigate("/profilo/" + booking.otherUser.uid)}
                                className="text-center font-bold"
                            >
                                {booking.otherUser.completeName}
                            </p>{" "}
                            <p
                                onClick={() =>
                                    booking.owner.uid !== booking.otherUser.uid
                                        ? navigate("/dashboard/" + booking.gift.id)
                                        : navigate("/prodotto/" + booking.gift.id)
                                }
                                className="hidden text-base underline md:flex"
                            >
                                {" "}
                                {booking.gift.name}
                            </p>
                        </div>
                        {!isMobile && (
                            <ProfilePic
                                w={14}
                                image={booking.otherUser.image}
                                userId={booking.otherUser.uid}
                                hiddenLvl={true}
                            />
                        )}
                    </div>
                </div>
                <div ref={fixedDivRef} className="z-10 w-full">
                    <div className="w-full md:hidden">
                        <BannerRegalo booking={booking} />
                    </div>
                    <div className="w-full bg-white pb-2 md:hidden">
                        {(auth.currentUser.uid === booking.receiver.uid || booking.isRequested) &&
                            !booking.isCanceled && <ConfirmBookingTile booking={booking} />}
                    </div>
                </div>
            </div>

            <div
                id="chatdiv"
                className={`flex flex-1 flex-col-reverse overflow-y-auto py-4 ${isMobile && "pt-64"}`}
                ref={scrollRef}
            >
                {loaded.current && messages.length > 0 ? (
                    messages.map((item: Message, i: number) => {
                        return (
                            //TODO: Enforce message ordering
                            //TODO: Optimize conditions on messages spacing and date display
                            //TODO: Messages date is compared only basing on month day
                            <>
                                <RenderMessage
                                    user={item.sendBy}
                                    message={item}
                                    gap={
                                        (i < messages.length - 1 &&
                                            item.createdAt.getDate() !==
                                                messages[i + 1].createdAt.getDate()) ||
                                        !(
                                            i < messages.length - 1 &&
                                            item.sendBy === messages[i + 1].sendBy
                                        )
                                    }
                                />
                                {i === messages.length - 1 ||
                                (i < messages.length - 1 &&
                                    item.createdAt.getDate() !==
                                        messages[i + 1].createdAt.getDate()) ? (
                                    <div className="flex justify-center">
                                        <div className="mt-8 select-none rounded-lg bg-grigino px-4 py-1 text-xs font-light text-gray-600">
                                            {item.createdAt.toLocaleString("default", {
                                                weekday: "long",
                                                day: "numeric",
                                                month: "long",
                                            })}
                                        </div>
                                    </div>
                                ) : (
                                    <div />
                                )}
                            </>
                        );
                    })
                ) : loaded.current && messages.length === 0 ? (
                    <div className="flex h-full w-full select-none items-center justify-center text-grigio">
                        Ancora nessun messaggio
                    </div>
                ) : (
                    <Caricamento transparent={true} />
                )}
            </div>
            {imagePreview && (
                <div className="flex justify-between px-4">
                    <img
                        src={imagePreview}
                        alt="Selected photo"
                        className="h-20 w-auto bg-blue-300"
                    />
                    <button
                        onClick={() => {
                            setImagePreview(null);
                            tmpImageFile = null;
                        }}
                    >
                        X
                    </button>
                </div>
            )}
            <div className="hidden w-full bg-white md:flex">
                {(auth.currentUser.uid === booking.receiver.uid || booking.isRequested) &&
                    !booking.isCanceled && <ConfirmBookingTile booking={booking} />}
            </div>
            <div
                className={`${isMobile && "fixed bottom-0 left-0 w-full"} h-16 rounded-lg bg-white ${booking.isBlocked ? "pointer-events-none opacity-50" : ""}`}
            >
                <form className={`flex items-center rounded-b-lg px-4 py-2`}>
                    <input
                        type="file"
                        accept={acceptedImageTypes.join(", ")}
                        ref={fileInputRef}
                        style={{ display: "none" }}
                        onChange={onPhotoChoose}
                    />
                    <button
                        className="cursor-pointer text-[#333]"
                        onClick={e => {
                            e.preventDefault();
                            if (fileInputRef.current) {
                                fileInputRef.current.click();
                            }
                        }}
                        type="button"
                    >
                        <Camera w={30} />
                    </button>
                    <input
                        className="border- mx-4 block w-full rounded-lg border-gray bg-white p-2.5 text-gray-900 focus:border-verde focus:ring-verde"
                        type="text"
                        value={text}
                        onChange={e => setText(e.target.value)}
                        id="chat"
                        placeholder="Scrivi qui il tuo messaggio..."
                        ref={textInputRef}
                    ></input>
                    <DefaultButton onClick={e => send(e)}>
                        <span className="flex items-center gap-2">
                            <p className="hidden md:flex">Invia</p>
                            <Send />
                        </span>
                    </DefaultButton>
                </form>
            </div>
        </div>
    );
}

export default ChatMessages;
