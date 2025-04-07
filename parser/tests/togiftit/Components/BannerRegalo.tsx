import React, { useEffect, useState } from "react";
import { Booking } from "../../../core/types/Exchange";
import GiftAPI from "../../../core/API/GiftAPI";
import { Gift } from "../../../core/types/Gift";
import Location from "../../../assets/Icons/Location";
import { useNavigate } from "react-router-dom";
import useUserStore from "../../../stores/userStore";

interface BannerRegaloProps {
    booking?: Booking;
    gift?: Gift;
}

const BannerRegalo: React.FC<BannerRegaloProps> = ({ booking = null, gift = null }) => {
    const [giftToShow, setGiftToShow] = useState<Gift>();
    const navigate = useNavigate();
    const userID = useUserStore(state => state.user?.uid);

    const handleLoad = async () => {
        const _regalo = await GiftAPI.getById(booking.gift.id);
        setGiftToShow(_regalo);
    };

    useEffect(() => {
        if (!gift) {
            handleLoad();
        } else {
            setGiftToShow(gift);
        }
    }, []);

    return (
        <div
            onClick={() =>
                navigate(
                    userID === giftToShow.owner
                        ? "/dashboard/" + giftToShow._id
                        : "/prodotto/" + giftToShow._id,
                )
            }
            className="w-full cursor-pointer bg-gray-ultralight py-2 md:rounded-lg md:bg-white md:px-5"
        >
            {giftToShow?._id && (
                <div className="flex h-20">
                    <img
                        //@ts-ignore
                        src={gift ? giftToShow.images[0].image : giftToShow?.images[0]}
                        alt=""
                        className="aspect-[4/3] rounded-lg object-cover"
                    />
                    <div className="ml-5 flex flex-col justify-center">
                        <p className="text-2xl font-bold">{giftToShow?.name}</p>
                        <span className="flex items-center gap-1">
                            <Location w={18} />
                            <p className="text-sm text-grigio">{giftToShow?.location.city}</p>
                        </span>
                    </div>
                </div>
            )}
        </div>
    );
};

export default BannerRegalo;
