import BookingAPI, { BookingStateUpdate } from "../../../core/API/BookingAPI";

export interface BookingActionButton {
    text: string;
    action?: (bookingId: string) => Promise<BookingStateUpdate>;
    successMessage?: string;
    errorMessage?: string;
}

export interface BookingActionTile {
    statusText: string;
    cancel?: BookingActionButton;
    confirm?: BookingActionButton;
}

export const genericBookingActions = {
    review: {
        statusText: "Scambio effettuato!",
        confirm: {
            text: "Lascia una recensione",
        },
    },
};

export const receiverBookingActions = {
    undoDelivery: {
        statusText: "Hai confermato il ritiro, in attesa della conferma di consegna.",
        cancel: {
            text: "Annulla conferma di ritiro",
            action: BookingAPI.undoGiftDelivery,
            successMessage: "Conferma di ritiro annullata",
            errorMessage: "Si è verificato un errore durante l'annullamento",
        },
    },
    confirmDelivery: {
        statusText: "Hai ritirato il regalo?",
        confirm: {
            text: "Conferma ritiro",
            action: BookingAPI.confirmGiftDelivery,
            successMessage: "Ritiro confermato!",
            errorMessage: "Si è verificato un errore durante la conferma",
        },
        cancel: {
            text: "Rinuncia assegnazione",
        },
    },
    cancel: {
        statusText: "La tua richiesta è in attesa di assegnazione.",
        cancel: {
            text: "Annulla prenotazione",
            action: BookingAPI.cancelBooking,
            successMessage: "Hai annullato la prenotazione",
            errorMessage: "Si è verificato un errore durante la cancellazione",
        },
    },
    book: {
        statusText: "Sei interessato a questo prodotto?",
        confirm: {
            text: "Prenotati",
            action: BookingAPI.updateInfoToBooking,
            successMessage: "Hai effettuato la prenotazione!",
            errorMessage: "Si è verificato un errore durante la prenotazione",
        },
    },
};

export const ownerBookingActions = {
    undoDelivery: {
        statusText: "Hai confermato la consegna del regalo, in attesa della conferma di ritiro.",
        cancel: {
            text: "Annulla conferma di consegna",
            action: BookingAPI.undoGiftDelivery,
            successMessage: "Conferma di consegna annullata",
            errorMessage: "Si è verificato un errore durante l'annullamento",
        },
    },
    confirmDelivery: {
        statusText: "Hai consegnato il regalo?",
        confirm: {
            text: "Conferma consegna",
            action: BookingAPI.confirmGiftDelivery,
            successMessage: "Consegna confermata!",
            errorMessage: "Si è verificato un errore durante la conferma",
        },
        cancel: {
            text: "Ritira Assegnazione",
        },
    },
    assign: {
        statusText: "Vuoi assegnare questo regalo a OTHER_USER_NAME? 💭",
        confirm: {
            text: "Assegna regalo",
            action: BookingAPI.assignGift,
            successMessage: "Regalo assegnato!",
            errorMessage: "Si è verificato un errore durante l'assegnazione",
        },
    },
};

export const parseBookingActionMessage = (message: string, otherUserName: string) => {
    return message.replace("OTHER_USER_NAME", otherUserName);
};
