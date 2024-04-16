import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import store from './store';
import {
    getProjList,
    getChatList,
    getMsgList,
    getNewMsg,
} from './comms';


function activate(context) {

    vscode.commands.executeCommand('gemify.openPanel');

    //

    let launch = vscode.commands.registerCommand('gemify.openPanel', function () {

        var gemify = vscode.window.createWebviewPanel(
            'gemify',
            'Gemify',
            vscode.ViewColumn.One,
            {
                enableScripts: true,
                retainContextWhenHidden: true,
            }
        );

        loadWebviewResources(context, gemify);

        gemify.webview.onDidReceiveMessage(msg => {
            switch (msg.command) {
                case 'execProjList':
                    getProjList()
                        .then(res => {
                            gemify.webview.postMessage({
                                command: 'returnProjList',
                                data: res.data,
                            });
                        })
                        .catch(err => {
                            console.error(err)
                        })
                    break;
                case 'execChatList':
                    var projID = msg.data.projID;
                    getChatList(projID)
                        .then(res => {
                            gemify.webview.postMessage({
                                command: 'returnChatList',
                                data: {
                                    chats: res.data.chats,
                                    proj: msg.data.proj,
                                }
                            });
                        })
                        .catch(err => {
                            console.error(err)
                        })
                    break;
                case 'execMsgList':
                    var projID = msg.data.proj.projID;
                    var chatID = msg.data.chat.chatID;
                    getMsgList(projID, chatID)
                        .then(res => {
                            gemify.webview.postMessage({
                                command: 'returnMsgList',
                                data: {
                                    msgs: res.data.messages,
                                    chat: msg.data.chat,
                                }
                            });
                        })
                        .catch(err => {
                            console.error(err)
                        })
                    break;
                case 'execNewMsg':
                    getNewMsg(msg, (chunk) => {
                        gemify.webview.postMessage({
                            command: 'returnMsg',
                            data: chunk,
                        });
                    });
                    break;


                    

                // case 'execHomeView':
                //     store.getState().setActiveConvoId(null)
                //     break;


                // case 'execReturnHome':
                //     store.getState().setActiveConvoId(null)
                //     break;

                // case 'execNewProj':
                //     gemify.webview.postMessage({
                //         command: 'returnNewProj',
                //     });
                //     break;



                // case 'execConvoView':
                //     const convoID = msg.data
                //     store.getState().setActiveConvoId(convoID)
                //     gemify.webview.postMessage({
                //         command: 'returnConvoView',
                //         data: convoID,
                //     });
                //     break;


                // case 'execHomeView':
                //     store.getState().setActiveConvoId(null)
                //     break;

                // case 'execNewConvo':
                //     sendNewConvo()
                //         .then(res => {
                //             const convoId = res.data.convoID;
                //             store.getState().setActiveConvoId(convoId)
                //         })
                //         .catch(err => {
                //             console.error(err)
                //         })
                //     break;



            }
        });
    });

    context.subscriptions.push(launch);
};

function loadWebviewResources(context, gemify) {

    const reset = gemify.webview.asWebviewUri(
        vscode.Uri.file(
            path.join(
                context.extensionPath,
                './src/styles/reset.css'
            )
        )
    );

    const css = gemify.webview.asWebviewUri(
        vscode.Uri.file(
            path.join(
                context.extensionPath,
                './src/styles/webview.css'
            )
        )
    );

    const js = gemify.webview.asWebviewUri(
        vscode.Uri.file(
            path.join(
                context.extensionPath,
                './src/webview.js'
            )
        )
    );

    const html = fs.readFileSync(
        vscode.Uri.file(
            path.join(
                context.extensionPath,
                './src/webview.html'
            )
        ).fsPath, 'utf-8');

    gemify.webview.html = /*html*/`<!DOCTYPE html>
    <html lang="en">
        <head>
            <link rel="stylesheet" href="${reset}" />
            <link rel="stylesheet" href="${css}" />
            <script src="${js}"></script>
            <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
            <title></title>
        </head>
        <body>
            ${html}
        </body>
    </html>`;
};

function deactivate() { }

module.exports = {
    activate,
    deactivate
};
