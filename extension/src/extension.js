import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import {
    sendGeminiReq
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
        )

        loadWebviewResources(context, gemify)

        // C2 - (Command and Control) 
        gemify.webview.onDidReceiveMessage(msg => {
            switch (msg.command) {
                case 'execGeminiMsg':
                    sendGeminiReq(msg.message)
                        .then(res => {
                            gemify.webview.postMessage({
                                command: 'displayGeminiRes',
                                data: {
                                    message: res
                                },
                                status: 'success'
                            });
                        })
                        .catch(err => {
                            console.error(err)
                        })
                    break;
            }
        });
    });

    context.subscriptions.push(launch);
}

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
            <title></title>
        </head>
        <body>
            ${html}
        </body>
    </html>`;
}

function deactivate() { }

module.exports = {
    activate,
    deactivate
}
