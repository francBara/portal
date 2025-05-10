import React from "react";
import BottoneRosso from "../Bottoni/BottoneRosso";
import Punti from "../../../assets/svg/Punti";

//@portal ui

function CardLanding({ props }) {
    const img = props.img[0].length
        ? props.img[0]
        : props.img[1].length
            ? props.img[1]
            : props.img[2].length
                ? props.img[2]
                : props.img[3].length
                    ? props.img[3]
                    : props.img[4].length
                        ? props.img[4]
                        : props.img[5]
                            ? props.img[5]
                            : require("../../../assets/default/default-image.jpg");

    return (
        <div>
            <div className="m-2 w-56 rounded-lg border border-gray-200 bg-white shadow-xl transition-all duration-200 md:w-[100%] md:hover:scale-105">
                <div
                    className="flex items-center justify-center p-4"
                    onClick={() => (window.location.href = "/prodotto/" + props._id)}
                >
                    <img
                        className="h-32 w-96 cursor-pointer rounded-lg object-cover md:h-52"
                        src={img}
                        alt={props.titolo}
                        loading="lazy"
                    />
                </div>
                <div className="cursor-pointer px-5 pb-5">
                    <h5
                        onClick={() => (window.location.href = "/prodotto/" + props._id)}
                        className="h-8 cursor-pointer text-xl font-semibold tracking-tight text-gray-900 dark:text-white"
                    >
                        {props.titolo.slice(0, 15)}
                        {props.titolo.length > 15 && "..."}
                    </h5>
                    <div
                        onClick={() => (window.location.href = "/prodotto/" + props._id)}
                        className="mb-3 mt-2.5 flex items-center space-x-1 font-semibold capitalize md:mb-5"
                    >
                        {props.categoria}
                    </div>
                    <div
                        onClick={() => (window.location.href = "/prodotto/" + props._id)}
                        className="mb-5 mt-2.5 flex h-1 items-center space-x-1 text-xs font-semibold capitalize"
                    >
                        {props.comune}, {props.provincia}
                    </div>
                    <div
                        onClick={() => (window.location.href = "/prodotto/" + props._id)}
                        className="flex items-center justify-between"
                    >
                        <span className="flex items-center space-x-2 text-2xl font-bold text-gray-900 dark:text-white">
                            <p>{props.price}</p> <Punti />
                        </span>

                        <BottoneRosso
                            text={"Apri"}
                            onclick={() => (window.location.href = "/prodotto/" + props._id)}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}

export default CardLanding;
