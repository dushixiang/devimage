import {Button, message, notification, Space, Table, Tag, Tooltip, Typography, Upload, UploadProps} from 'antd';
import React, {useState} from 'react';
import {
    CheckCircleOutlined,
    ClearOutlined,
    CloseCircleOutlined,
    InboxOutlined,
    SearchOutlined,
    SyncOutlined
} from "@ant-design/icons";
import type {RcFile} from 'antd/es/upload/interface';
import {AddFile, ClearCompressed, OpenFileDialog} from "../../wailsjs/go/main/App";
import {EventsOff, EventsOn, LogPrint} from "../../wailsjs/runtime";

type FileProp = {
    file: string,
    size: number,
    compressed: number,
    ratio: number,
    status: string
    error: string
}

const renderSize = (value: number) => {
    if (null == value || value === 0) {
        return "0 B";
    }
    const unitArr = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
    let index = Math.floor(Math.log(value) / Math.log(1024));
    let size = value / Math.pow(1024, index);
    let sizeStr = size.toFixed(2);
    return sizeStr + ' ' + unitArr[index];
}

const Compress = () => {

    const [api, contextHolder] = notification.useNotification();
    let [files, setFiles] = useState<FileProp[]>([]);

    const openNotification = (message: string) => {
        api.success({
            message: `Notification`,
            description: message,
            placement: 'bottomRight',
        });
    };

    const toBase64 = (file: any) => new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });

    const beforeUpload = async (file: RcFile) => {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
        if (!isJpgOrPng) {
            message.warning('You can only upload JPG/PNG file!');
        } else {
            files.unshift({
                compressed: 0, error: "", file: file.name, ratio: 0, size: file.size, status: "processing"
            })
            LogPrint(JSON.stringify(files));
            setFiles(files.slice());
            let raw = await toBase64(file);
            let fileInfo = {
                file: file.name,
                size: file.size,
                raw: raw,
            }
            let s = JSON.stringify(fileInfo);
            AddFile(s).then((result) => {
                if (result['code'] !== 200) {
                    message.error(result['message']);
                    return {};
                }
            })
        }
        return false;
    };

    const props: UploadProps = {
        name: 'file',
        multiple: true,
        beforeUpload: beforeUpload,
        openFileDialogOnClick: false,
        itemRender: () => false,
        onDrop(e) {
            console.log('Dropped files', e.dataTransfer.files);
            reListenEvent()
        },
    };

    const columns = [
        {
            title: 'FILE',
            dataIndex: 'file',
            key: 'file',
            ellipsis: true,
        },
        {
            title: 'SIZE',
            dataIndex: 'size',
            key: 'size',
            render: (text: number) => {
                return renderSize(text);
            }
        },
        {
            title: 'COMPRESSED',
            dataIndex: 'compressed',
            key: 'compressed',
            render: (text: number) => {
                return renderSize(text);
            }
        },
        {
            title: 'RATIO',
            dataIndex: 'ratio',
            key: 'ratio',
            render: (text: number) => {
                return text + '%';
            }
        },
        {
            title: 'STATUS',
            dataIndex: 'status',
            key: 'status',
            render: (text: string, record: FileProp) => {
                let display = <div>{text}</div>
                switch (text) {
                    case 'success':
                        display = <Tag key={'success'} icon={<CheckCircleOutlined/>} color="success">
                            Success
                        </Tag>;
                        break;
                    case 'processing':
                        display = <Tag key={'processing'} icon={<SyncOutlined spin/>} color="processing">
                            Processing
                        </Tag>;
                        break;
                    case 'error':
                        display = <Tag key={'error'} icon={<CloseCircleOutlined/>} color="error">
                            Error
                        </Tag>
                        break;
                }
                return <Tooltip title={record.error}>{display}</Tooltip>
            }
        },
    ];

    function reListenEvent() {
        EventsOff("CompressChange")
        EventsOn("CompressChange", (data: FileProp[]) => {
            setFiles(data);
        })
    }

    return (
        <div>
            {contextHolder}
            <div style={{display: "flex"}}>
                <div style={{display: "flex", flex: '1 1 0%'}}>
                    <Typography.Title level={3}>Compress</Typography.Title>
                </div>
                <div style={{textAlign: 'left'}}>
                    <Space>
                        <Button type={'text'}
                                icon={<ClearOutlined/>}
                                onClick={() => {
                                    ClearCompressed().then(r => openNotification('clear success.'))
                                }}
                        >

                        </Button>
                        {/*<Button type={'text'} icon={<SearchOutlined/>}>*/}

                        {/*</Button>*/}
                    </Space>
                </div>
            </div>


            <div onClick={(event) => {
                reListenEvent();
                OpenFileDialog().then(r => {
                    setFiles(r['data'] as FileProp[]);
                })
                return false;
            }}>
                <Upload.Dragger {...props} style={{height: 250}}>
                    <p className="ant-upload-drag-icon">
                        <InboxOutlined/>
                    </p>
                    <p className="ant-upload-text">Click or drag file to this area to upload</p>
                    <p className="ant-upload-hint">
                        Support for a single or bulk upload. Strictly prohibit from uploading company data or other
                        band files
                    </p>
                </Upload.Dragger>
            </div>

            <div style={{marginTop: 16}}>
                <Table size={'small'}
                       dataSource={files}
                       columns={columns}
                       pagination={false}
                       rowKey={'file'}
                       scroll={{y: 'calc(100vh - 280px)'}}
                />
            </div>
        </div>
    );
};

export default Compress;