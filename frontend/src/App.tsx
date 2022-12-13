import React from 'react';
import './App.css';
import {Button, Layout, Row, Space, theme} from "antd";
import {InfoCircleOutlined, OneToOneOutlined, SettingOutlined} from "@ant-design/icons";
import {Link, Outlet} from "react-router-dom";

const isMac = () => {
    return /macintosh|mac os x/i.test(navigator.userAgent);
}

function App() {
    const {
        token: {colorBgContainer},
    } = theme.useToken();

    let marginTop = 0;
    if(isMac()){
        marginTop = 36;
    }

    return (
        <>
            <Layout style={{background: colorBgContainer, height: '100vh'}}>
                <Layout.Sider style={{background: colorBgContainer}} width={50}>
                    <Row justify="center" align="top">
                        <Space direction={'vertical'} size={'small'} style={{marginTop: marginTop}}>
                            <Link to={'/setting'}>
                                <Button type={'text'} icon={<SettingOutlined/>}/>
                            </Link>
                            <Link to={'/'}>
                                <Button type={'text'} icon={<OneToOneOutlined/>}/>
                            </Link>
                            <Link to={'/about'}>
                                <Button type={'text'} icon={<InfoCircleOutlined/>}/>
                            </Link>
                        </Space>
                    </Row>
                </Layout.Sider>
                <Layout.Content style={{margin: 16, marginTop: marginTop}}>
                    <Outlet/>
                </Layout.Content>
            </Layout>
        </>
    );
}

export default App
