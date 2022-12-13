import React from 'react';
import {Descriptions, Typography} from "antd";

const About = () => {
    return (
        <div>
            <Typography.Title level={3}>About</Typography.Title>
            <Descriptions column={1}>
                <Descriptions.Item label="Version">v0.0.1</Descriptions.Item>
                <Descriptions.Item label="Author">dushixiang</Descriptions.Item>
                <Descriptions.Item label="Github">https://github.com/dushixiang</Descriptions.Item>
                <Descriptions.Item label="Website">https://typesafe.cn</Descriptions.Item>
            </Descriptions>
        </div>
    );
};

export default About;