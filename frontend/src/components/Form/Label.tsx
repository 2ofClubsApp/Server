import React from 'react'
import Form from "react-bootstrap/Form";
import {FormInfo} from "../../types/FormInfo";

type FormProps = {
    key: number
    props: FormInfo
}

const Label = (label: FormProps) => {
    return (
        <Form.Group controlId={label.props.controlId}>
            <Form.Label>{label.props.label}</Form.Label>
            <Form.Control type={label.props.type} placeholder={label.props.placeholder}/>
        </Form.Group>
    )
};

export default Label;
