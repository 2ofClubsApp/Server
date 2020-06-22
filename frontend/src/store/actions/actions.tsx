import {SET_USERDATA} from "./actionTypes";
import React from "react";

export const setUserData = (id: string, value: string) => {
    return (
        {
            type: SET_USERDATA,
            payload: {id: value}
        }
    )
}

