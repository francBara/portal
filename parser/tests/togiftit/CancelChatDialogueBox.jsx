import React from "react";
import { useState } from "react";
import useDeviceDetection from "../../core/useDeviceDetection";
import DragCloseDrawer from "../../Components/Atoms/Drawers/DragCloseDrawer";

function ReasonTile({ text, isSelected, onClick }) {
    return (
        <div
            className="flex cursor-pointer flex-row items-center justify-between"
            onClick={onClick}
        >
            <div className="text-md mr-4">{text}</div>
            <input
                type="checkbox"
                name=""
                checked={isSelected}
                className="appearance-none checked:border-0 checked:bg-accent checked:ring-0 focus:ring-0"
                id=""
            />
        </div>
    );
}

function CancelChatDialogueBox({ setIsDialogueVisible, isReceiver, onConfirm }) {
    const [reason, setReason] = useState(null);
    const [customReason, setCustomReason] = useState("");
    const [isConfirmClickable, setConfirmClickable] = useState(false);
    const device = useDeviceDetection();

    const reasons = isReceiver
        ? [
              "Ho avuto un contrattempo",
              "Non sono più interessato al regalo",
              "Il regalo non è stato presentato chiaramente",
          ]
        : ["Ho avuto un contrattempo", "Il ricevitore non ha dato istruzioni chiare"];
    /*
    const checkConfirmClickable = () => {
        setConfirmClickable(reason !== null || customReason.length > 0);
    };
*/
    const content = (
        <div className="px-4">
            <div className="text-xl font-semibold">
                {isReceiver
                    ? "Perché non vuoi ritirare il tuo regalo?"
                    : "Perché non vuoi consegnare il tuo regalo?"}
            </div>
            <div className="mt-4 flex w-full flex-col gap-3">
                {reasons.map((reasonText, index) => {
                    return (
                        <ReasonTile
                            text={reasonText}
                            isSelected={reason === index}
                            onClick={() => {
                                if (index === reason) {
                                    setReason(null);
                                    if (customReason.length === 0) {
                                        setConfirmClickable(false);
                                    }
                                } else {
                                    setReason(index);
                                    setConfirmClickable(true);
                                }
                            }}
                        />
                    );
                })}
            </div>
            <div className="mt-2 flex w-full flex-col">
                <text className="font-medium">Altro</text>
                <input
                    className="mt-2 rounded-lg border-0 bg-gray-100 p-2.5 text-sm text-gray-900"
                    type="text"
                    value={customReason}
                    onChange={e => {
                        setCustomReason(e.target.value);
                        if (e.target.value.length > 0) {
                            setConfirmClickable(true);
                        } else if (reason === null) {
                            setConfirmClickable(false);
                        }
                    }}
                    id="chat"
                    placeholder="Scrivi qui il tuo messaggio..."
                ></input>
            </div>
            <div className="flex w-full justify-center">
                <button
                    className={`mt-8 ${
                        isConfirmClickable
                            ? "bg-accent text-white"
                            : "cursor-default bg-grigino text-black text-opacity-50"
                    } w-full rounded-md px-4 py-2 sm:w-auto`}
                    onClick={() => {
                        if (customReason.length > 0) {
                            onConfirm(customReason);
                        } else if (reason !== null) {
                            onConfirm(reasons[reason]);
                        }
                    }}
                >
                    Conferma
                </button>
            </div>
        </div>
    );

    if (device === "Mobile")
        return (
            <DragCloseDrawer open={true} setOpen={setIsDialogueVisible}>
                {content}
            </DragCloseDrawer>
        );
    return (
        <div
            className="fixed inset-0 z-50 grid place-items-center backdrop-blur-sm"
            onClick={() => {
                setIsDialogueVisible(false);
            }}
        >
            <div
                className="flex w-[95vw] flex-col items-start rounded-md border-2 bg-white p-4 md:w-auto"
                onClick={e => e.stopPropagation()}
            >
                {content}
            </div>
        </div>
    );
}

export default CancelChatDialogueBox;
