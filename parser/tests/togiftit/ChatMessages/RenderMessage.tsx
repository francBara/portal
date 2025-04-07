import React, { useState } from "react";
import { auth } from "../../../firebase";
import FadeInComponent from "../../../Components/Atoms/Transitions/FadeInComponent";
import { renderTime } from "../../../core/StringUtils";
import { AnimatePresence, motion } from "framer-motion";
import { IoIosClose } from "react-icons/io";
import ProfilePic from "../../../Components/Atoms/ProfilePic/ProfilePic";
import { Message } from "../../../core/types/Chat";
import { User } from "../../../core/types/User";
/**
 * @param {Object}o props
 * @param {import('../../../core/API/ChatAPI').Message} props.messager
 * @param {boolean} props.isRight
 * @returns
 */

interface ChatMessageProps {
    message: Message & { isLoading?: boolean };
    isRight: boolean;
    user?: User;
}

function ChatMessage({ message, isRight, user = null }: ChatMessageProps) {
    const [imageDisplayed, setImageDisplayed] = useState(false);
    const [isExpanded, setIsExpanded] = useState(false);

    const toggleExpand = () => {
        setIsExpanded(!isExpanded);
    };

    const truncatedContent =
        message.content.length > 250 ? message.content.substring(0, 250) + "..." : message.content;
    console.log(user);

    return (
        <div className="flex flex-col">
            <div className="flex items-end gap-2">
                {!isRight &&
                    (user ? (
                        <span className="relative z-[0] flex items-end">
                            <ProfilePic w={8} hiddenLvl hiddenLH userId={user} />
                        </span>
                    ) : (
                        <span className="relative z-[0] w-8" />
                    ))}
                <div
                    className={`flex max-w-[80vw] flex-col rounded-lg ${isRight ? "bg-gray-light text-black" : "bg-secondary text-white"} ${message.isLoading ? "opacity-40" : "opacity-100"}`}
                >
                    {message.image ? (
                        <div className="p-2">
                            <img
                                src={message.image}
                                alt="Sent image"
                                className={`h-20 w-auto cursor-pointer transition duration-300 hover:brightness-125 ${message.content.length > 0 ? "rounded-md" : "rounded-md"} ${message.isLoading ? "opacity-40" : "opacity-100"}`}
                                onClick={() => {
                                    setImageDisplayed(true);
                                }}
                            />
                            <ChatDisplayImage
                                image={message.image}
                                isOpen={imageDisplayed}
                                setIsOpen={setImageDisplayed}
                            />
                        </div>
                    ) : (
                        <div />
                    )}
                    <div className="flex items-end justify-between overflow-hidden px-3 py-2">
                        {message.content.length > 0 ? (
                            <div className="max-w-48 break-words pr-2 text-left text-base font-normal leading-5 sm:max-w-sm">
                                {isExpanded ? message.content : truncatedContent}
                                {message.content.length > 250 && (
                                    <span
                                        className="cursor-pointer text-blue-500"
                                        onClick={toggleExpand}
                                    >
                                        {isExpanded ? " Mostra meno" : " Mostra di più"}
                                    </span>
                                )}
                            </div>
                        ) : (
                            <div />
                        )}
                    </div>
                </div>
            </div>
            <div
                className={`text-xs ${isRight ? "self-end" : "self-start pl-10"} mt-1 font-extralight text-gray-500`}
            >
                {renderTime(message.createdAt)}
            </div>
        </div>
    );
}

function ChatDisplayImage({ isOpen, setIsOpen, image }) {
    return (
        <AnimatePresence>
            (
            {isOpen && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    onClick={() => setIsOpen(false)}
                    className="fixed inset-0 z-50 grid cursor-pointer place-items-center backdrop-blur"
                >
                    <motion.div
                        initial={{ scale: 0, rotate: "0deg" }}
                        animate={{ scale: 1, rotate: "0deg" }}
                        exit={{ scale: 0, rotate: "0deg" }}
                        onClick={e => e.stopPropagation()}
                        className="relative flex h-full w-full cursor-default items-center justify-center overflow-hidden bg-white p-2 text-white shadow-xl"
                    >
                        <IoIosClose
                            onClick={() => setIsOpen(false)}
                            className="absolute right-3 top-3 h-12 w-12 cursor-pointer rounded-full bg-transparent fill-black hover:bg-accent"
                        />
                        <div className="">
                            <img
                                src={image}
                                className="h-full w-full object-contain"
                                alt="Immagine non caricata correttamente"
                            />
                        </div>
                    </motion.div>
                </motion.div>
            )}
            )
        </AnimatePresence>
    );
}

/**
 * @param {Object} props
 * @param {import('../../../core/API/ChatAPI').Message} props.message
 */
function AutoMessage({ message }) {
    const autoTypes = {
        confirmedOwner: {
            title: "Consegna confermata! 🌱",
            text: "Il donatore ha consegnato il regalo",
        },
        confirmedReceiver: { title: "Ritiro confermato! 🌱", text: "Il regalo è stato ricevuto" },
        canceled: { title: "Prenotazione annullata", text: "La prenotazione è stata annullata" },
        deliveryUndone: {
            title: "Conferma di consegna annullata",
            text: "La conferma di consegna è stata annullata",
        },
        booked: {
            title: "Prenotazione effettuata! 🙋‍♀️",
            text: "Utilizzate questa chat per accordarvi",
        },
        assigned: {
            title: "Assegnazione effettuata! 😊",
            text: "Una volta fatto, confermate lo scambio",
        },
        info: { title: "Chiedi informazioni sul regalo! 💬", text: "" },
    };

    const autoType = autoTypes[message.autoType];

    return (
        <div className="mb-8 flex flex-col items-center">
            <div className="text-sm font-semibold">{autoType.title}</div>
            <div className="text-sm">{autoType.text}</div>
        </div>
    );
}

function RenderMessage({ message, gap, user }) {
    return (
        <FadeInComponent>
            <div className={`mt-${gap ? 4 : 2} px-4`}>
                {message.autoType ? (
                    <div className="flex w-full justify-center">
                        <AutoMessage message={message} />
                    </div>
                ) : message.sendBy === auth.currentUser.uid ? (
                    <div className="flex w-full flex-row-reverse">
                        <ChatMessage message={message} isRight={true} user={user} />
                    </div>
                ) : (
                    <div className="flex overflow-hidden">
                        <ChatMessage message={message} isRight={false} user={gap && user} />
                    </div>
                )}
            </div>
        </FadeInComponent>
    );
}

export default RenderMessage;
