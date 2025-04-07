import { auth } from "../../../firebase";
import { useEffect, useState } from "react";
import CancelChatDialogueBox from "../CancelChatDialogueBox";
import React from "react";
import { ReviewSpringModal } from "../../../Components/Atoms/Modals/Recensisci";
import BookingAPI from "../../../core/API/BookingAPI";
import { toast } from "sonner";
import useBookingStore from "../../../stores/bookingStore";
import { showPrompt } from "../../../Components/Prompts/Prompts";
import AreYouSure from "../../../Components/Atoms/Modals/AreYouSure";
import {
    BookingActionButton,
    BookingActionTile,
    genericBookingActions,
    ownerBookingActions,
    parseBookingActionMessage,
    receiverBookingActions,
} from "./BookingActions";
import { Booking } from "../../../core/types/Exchange";
import useDeviceDetection from "../../../core/useDeviceDetection";

interface ConfirmBookingTileProps {
    booking: Booking;
}

function ConfirmBookingTile({ booking }: ConfirmBookingTileProps) {
    const [isDialogueVisible, setIsDialogueVisible] = useState(false);
    const [isReviewOpen, setIsReviewOpen] = useState(false);

    const isReceiver = booking.receiver.uid === auth.currentUser?.uid;

    const appendMessage = useBookingStore(state => state.appendMessage);
    const updateBookings = useBookingStore(state => state.updateBookings);

    const [statusText, setStatusText] = useState("");

    const [confirm, setConfirm] = useState<{
        text: string;
        action: () => void;
    } | null>(null);

    const [cancel, setCancel] = useState<{
        text: string;
        action: () => void;
    } | null>(null);

    const parseBookingAction = (actionButton: BookingActionButton) => {
        return async () => {
            try {
                const stateUpdate = await actionButton.action!(booking._id);

                if (stateUpdate.booking.isDelivered !== undefined) {
                    if (isReceiver) {
                        stateUpdate.booking.isConfirmedByReceiver = stateUpdate.booking.isDelivered;
                        stateUpdate.message.autoType = stateUpdate.booking.isDelivered
                            ? "confirmedReceiver"
                            : "deliveryUndone";
                    } else {
                        stateUpdate.booking.isConfirmedByOwner = stateUpdate.booking.isDelivered;
                        stateUpdate.message.autoType = stateUpdate.booking.isDelivered
                            ? "confirmedOwner"
                            : "deliveryUndone";
                    }
                }

                appendMessage(stateUpdate.message);
                const updatedBooking = { ...booking, ...stateUpdate.booking };

                updateBookings(updatedBooking);

                toast.success(
                    parseBookingActionMessage(
                        actionButton.successMessage!,
                        booking.otherUser.completeName,
                    ),
                );
                if (
                    stateUpdate.booking.isDelivered !== undefined &&
                    booking.isConfirmedByOwner &&
                    booking.isConfirmedByReceiver &&
                    !isReceiver
                ) {
                    showPrompt.pointsReceived();
                } else if (stateUpdate.booking.isRequested) {
                    showPrompt.pointsSpent();
                }
            } catch (e) {
                toast.error(
                    parseBookingActionMessage(
                        actionButton.errorMessage!,
                        booking.otherUser.completeName,
                    ),
                );
            }
        };
    };

    const updateTileState = (
        bookingAction: BookingActionTile,
        customActions?: { confirm?: () => void; cancel?: () => void },
    ) => {
        setStatusText(
            parseBookingActionMessage(bookingAction.statusText, booking.otherUser.completeName),
        );

        if (bookingAction.confirm) {
            if (bookingAction.confirm.action) {
                setConfirm({
                    text: bookingAction.confirm.text,
                    action: parseBookingAction(bookingAction.confirm),
                });
            } else if (customActions && customActions.confirm) {
                setConfirm({
                    text: bookingAction.confirm.text,
                    action: customActions.confirm,
                });
            }
        }
        if (bookingAction.cancel) {
            if (bookingAction.cancel.action) {
                setCancel({
                    text: bookingAction.cancel.text,
                    action: parseBookingAction(bookingAction.cancel),
                });
            } else if (customActions && customActions.cancel) {
                setCancel({
                    text: bookingAction.cancel.text,
                    action: customActions.cancel,
                });
            }
        }
    };

    const initTexts = () => {
        setStatusText("");
        setConfirm(null);
        setCancel(null);

        if (booking.isBlocked) {
            setStatusText("Non è possibile mandare messaggi");
            return;
        }

        if (booking.isConfirmedByOwner && booking.isConfirmedByReceiver) {
            setStatusText("Scambio effettuato!");
            if (
                (isReceiver && !booking.isReviewedByReceiver) ||
                (!isReceiver && !booking.isReviewedByOwner)
            ) {
                setConfirm({
                    text: "Lascia una recensione",
                    action: async () => {
                        setIsReviewOpen(true);
                    },
                });
            }
            return;
        }

        if (isReceiver) {
            if (booking.isAccepted) {
                if (booking.isConfirmedByReceiver) {
                    updateTileState(receiverBookingActions.undoDelivery);
                } else {
                    updateTileState(receiverBookingActions.confirmDelivery, {
                        cancel: async () => {
                            try {
                                const stateUpdate = await BookingAPI.cancelBooking(booking._id);
                                appendMessage(stateUpdate.message);
                                booking = { ...booking, ...stateUpdate.booking };
                                updateBookings(booking);
                                toast.success("Hai annullato la prenotazione");
                                setIsDialogueVisible(false);
                            } catch (e) {
                                toast.error("Si è verificato un errore durante la cancellazione");
                            }
                        },
                    });
                }
            } else if (booking.isRequested) {
                updateTileState(receiverBookingActions.cancel);
            } else {
                updateTileState(receiverBookingActions.book);
            }
        } else {
            if (booking.isAccepted) {
                if (booking.isConfirmedByOwner) {
                    updateTileState(ownerBookingActions.undoDelivery);
                } else {
                    updateTileState(ownerBookingActions.confirmDelivery, {
                        cancel: async () => {
                            try {
                                const stateUpdate = await BookingAPI.cancelBooking(booking._id);
                                appendMessage(stateUpdate.message);
                                booking = { ...booking, ...stateUpdate.booking };
                                updateBookings(booking);
                                toast.success("Hai annullato la prenotazione");
                                setIsDialogueVisible(false);
                            } catch (e) {
                                toast.error("Si è verificato un errore durante la cancellazione");
                            }
                        },
                    });
                }
            } else if (booking.isRequested) {
                updateTileState(ownerBookingActions.assign);
            }
        }
    };

    useEffect(() => {
        initTexts();
    }, [booking]);

    const isMobile = useDeviceDetection() === "Mobile";

    return (
        <div className="mt-0 flex w-full flex-col items-center justify-between bg-accent/20 px-4 py-4 md:flex-row">
            <ReviewSpringModal
                userUid={booking.otherUser.uid}
                userName={booking.otherUser.completeName}
                isOpen={isReviewOpen}
                setIsOpen={setIsReviewOpen}
                onReviewed={async () => {
                    await BookingAPI.setReviewed(booking._id);
                    if (isReceiver) {
                        booking.isReviewedByReceiver = true;
                    } else {
                        booking.isReviewedByOwner = true;
                    }
                    initTexts();
                    window.location.href = "/profilo/" + booking.otherUser.uid;
                }}
            />

            <div className="text-sm font-normal md:text-lg md:font-semibold">{statusText}</div>
            <div className="mt-2 flex w-full flex-col items-center gap-2 md:mt-0 md:w-auto md:flex-row md:gap-5">
                {cancel !== null && (
                    <div className="w-full md:w-auto">
                        <AreYouSure
                            dim={isMobile ? "small" : "normal"}
                            variant="red-secondary"
                            cancel
                            title={cancel.text}
                            head={cancel.text}
                            onClick={cancel.action}
                            text={statusText}
                        />
                    </div>
                )}
                {confirm !== null && (
                    <div className="w-full md:w-auto">
                        <AreYouSure
                            dim={isMobile ? "small" : "normal"}
                            withoutConfirm={confirm.text === "Lascia una recensione"}
                            title={confirm.text}
                            head={confirm.text}
                            onClick={confirm.action}
                            text={statusText}
                            variant={"green"}
                        />
                    </div>
                )}
            </div>

            {isDialogueVisible && (
                <CancelChatDialogueBox
                    setIsDialogueVisible={setIsDialogueVisible}
                    isReceiver={isReceiver}
                    onConfirm={async (reason: string) => {
                        try {
                            const stateUpdate = await BookingAPI.cancelBooking(booking._id, reason);
                            appendMessage(stateUpdate.message);
                            booking = { ...booking, ...stateUpdate.booking };
                            updateBookings(booking);
                            toast.success("Hai annullato la prenotazione");
                            setIsDialogueVisible(false);
                        } catch (e) {
                            toast.error("Si è verificato un errore durante la cancellazione");
                        }
                    }}
                />
            )}
        </div>
    );
}

export default ConfirmBookingTile;
