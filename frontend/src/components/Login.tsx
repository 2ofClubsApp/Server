import React from 'react'
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Formz from "./Signup/Form/Form";
import axios from '../axios'

const Login = () => {
    const propz = {
        controlId: "email",
        label: "Email Address",
        type: "email",
        placeholder: "Email"
    };
    return (
        <div>
            <Formz props={propz}/>
            <Form>
                <Form.Group controlId={"Email"}>
                    <Form.Label>Email Address</Form.Label>
                    <Form.Control type={"email"} placeholder={"Email"}/>
                </Form.Group>
            </Form>
            <Form>
                <Form.Group controlId={"Password"}>
                    <Form.Label>Password</Form.Label>
                    <Form.Control type={"password"} placeholder={"Password"}/>
                </Form.Group>
            </Form>
            <Button onClick={login}>Login</Button>
        </div>
    )
};

const login = () => {
    axios.get("/students/Chris").then(response => {
        console.log(response)
    });
    // console.log("Logging in")
};

export default Login
