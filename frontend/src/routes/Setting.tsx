import React from 'react';
import {Button, Form, Input, message, Spin, Typography} from "antd";
import {useQuery} from "react-query";
import {GetConfig, SetConfig} from "../../wailsjs/go/main/App";

const Setting = () => {

    const [form] = Form.useForm();

    let query = useQuery('get', async () => {
        let result = await GetConfig();
        if (result['code'] !== 200) {
            message.error(result['message']);
            return {};
        }
        return result['data'];
    }, {
        onSuccess: data => {
            form.setFieldsValue(data);
        }
    });

    const handleSave = (values: any) => {
        SetConfig(values).then(result => {
            if (result['code'] !== 200) {
                message.error(result['message']);
            }else {
                message.success(result['message']);
            }
        })
    }

    return (
        <div>
            <Typography.Title level={3}>Setting</Typography.Title>
            <Spin spinning={query.isLoading}>

                <Form form={form} name="setting" onFinish={handleSave} layout="vertical">
                    {/*<Form.Item*/}
                    {/*    name="quality"*/}
                    {/*    label="Quality"*/}
                    {/*    rules={[*/}
                    {/*        {*/}
                    {/*            required: true,*/}
                    {/*        },*/}
                    {/*    ]}*/}
                    {/*>*/}
                    {/*    <Slider*/}
                    {/*        min={1}*/}
                    {/*        max={100}*/}
                    {/*    />*/}
                    {/*</Form.Item>*/}

                    <Form.Item
                        name="outputDir"
                        label="OutputDir"
                        rules={[
                            {
                                required: true,
                            },
                        ]}
                    >
                        <Input/>
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            Save
                        </Button>
                    </Form.Item>
                </Form>
            </Spin>
        </div>
    );
};

export default Setting;