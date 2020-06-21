import {SET_USERNAME, SET_PASSWORD, SET_EMAIL} from "./actionTypes";

export const setUsername = (username: string) => ({
    type: SET_USERNAME,
    payload: {
        username
    }
})

export const setPassword = (password: string) => ({
    type: SET_USERNAME,
    payload: {
        password
    }
})

export const setEmail = (email: string) => ({
    type: SET_USERNAME,
    payload: {
        email
    }
})
