import React from 'react'
import Form from "react-bootstrap/Form";
import {FormInfo} from "../../types/FormInfo";

type FormProps = {
    props: FormInfo
}

const Formz = (a: FormProps) => {
    return (
        <Form>
            <Form.Group controlId={a.props.controlId}>
                <Form.Label>{a.props.label}</Form.Label>
                <Form.Control type={a.props.type} placeholder={a.props.placeholder}/>
            </Form.Group>
        </Form>
    )
};

export default Formz;
