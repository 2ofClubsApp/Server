import React from 'react'
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Button from "react-bootstrap/Button";
import Formz from "./Signup/Form/Form";
import axios from '../axios'
import {useHistory} from 'react-router-dom'
import '../app.css';


const Login = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };
    const email = {
        controlId: "email",
        label: "Email Address",
        type: "email",
        placeholder: "Email"
    };
    const pass = {
        controlId: "password",
        label: "Password",
        type: "password",
        placeholder: "Password"
    };
    return (
        <div> 
            <Button variant="outline-light" className="m-2 text-uppercase" onClick={() => changeRoute('/')}>Back to Home</Button>
            <Container className="form d-flex justify-content-center align-items-center" style={{width: "350px"}}>
                <Row>
                    <Col xs lg="12">
                        <h2 className="text-center mb-4">Login</h2>
                        <Formz props={email}/>
                        <Formz props={pass}/>
                        <div className="d-flex justify-content-center">
                        <Button onClick={login} className=" align-middle mt-3 pl-4 pr-4 pt-2 pd-2 btn-grad">Login</Button>
                        </div>
                    </Col>
                </Row>
            </Container>
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
