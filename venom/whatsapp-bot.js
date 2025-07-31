const venom = require('venom-bot');
const axios = require('axios');

const GO_API_URL = process.env.GO_API_URL || 'http://localhost:8080';
const BOT_SESSION = process.env.BOT_SESSION || 'mestra-bot';

let client = null;

async function sendQRCodeToAPI(base64Qr, attempts, urlCode) {
    try {
        console.log('QR Code gerado! Escaneie com seu WhatsApp.');
        
        await sendLog('info', 'QR Code gerado', {
            attempts,
            urlCode: urlCode?.substring(0, 50) + '...'
        });
        
    } catch (error) {
        console.error('Erro ao processar QR Code:', error.message);
        await sendLog('error', 'Erro ao processar QR Code', {
            error: error.message,
            attempts
        });
    }
}

async function updateQRCodeStatus(status) {
    try {
        const qrResponse = await axios.get(`${GO_API_URL}/api/v1/qr-code/latest`);
        if (qrResponse.data && qrResponse.data.qr_code) {
            const qrId = qrResponse.data.qr_code.id;
            
            await axios.put(`${GO_API_URL}/api/v1/qr-code/status`, {
                id: qrId,
                status: status
            });
            
            console.log(`Status do QR Code atualizado para: ${status}`);
        }
    } catch (error) {
        console.error('Erro ao atualizar status do QR Code:', error.message);
    }
}

async function sendLog(level, message, data = null) {
    try {
        await axios.post(`${GO_API_URL}/api/bot/log`, {
            level,
            message,
            data,
            timestamp: new Date().toISOString()
        });
    } catch (error) {
        console.error('Erro ao enviar log:', error.message);
    }
}

async function processMessage(message) {
    try {
        console.log('Mensagem recebida:', message);
        
        if (message.type === 'image' || message.mimetype?.startsWith('image/') || message.mediaData?.mimetype?.startsWith('image/')) {
            await processImageMessage(message);
            return;
        }
        
        if (message.type !== 'chat' || !message.body || message.body.trim() === '') {
            return;
        }
        
        console.log('Enviando para API:', {
            mensagem: message.body,
            autor: message.author || message.from,
            data: new Date().toLocaleDateString('pt-BR')
        });
        
        const response = await axios.post(`${GO_API_URL}/api/processar`, {
            mensagem: message.body,
            autor: message.author || message.from,
            data: new Date().toLocaleDateString('pt-BR')
        });

        console.log('Status HTTP:', response.status);
        console.log('Resposta da API:', response.data);
        console.log('Tipo da resposta:', typeof response.data);
        console.log('É array?', Array.isArray(response.data));

        // Validar se a API extraiu algum dado econômico
        if (response.data && Array.isArray(response.data) && response.data.length > 0) {
            // Sucesso - montar mensagem de confirmação com detalhes
            const itensExtracted = response.data.length;
            const totalValue = response.data.reduce((sum, item) => sum + item.valor, 0);
            
            // Montar lista de gastos validados
            let gastosLista = '';
            response.data.forEach((item, index) => {
                gastosLista += `${index + 1}. R$ ${item.valor.toFixed(2)} - ${item.descricao}\n`;
            });
            
            const confirmMsg = `✅ *${itensExtracted} gastos registrados com sucesso!*\n\n` +
                `� *Gastos validados:*\n${gastosLista}\n` +
                `💰 *Total: R$ ${totalValue.toFixed(2)}*\n\n` +
                `📝 *Dados salvos automaticamente!*`;
            
            await client.sendText(message.from, confirmMsg);
            
        } else {
            // Erro - nenhum dado econômico encontrado
            const errorMsg = `❌ *Ops! Não consegui extrair dados econômicos*\n\n` +
                `🔍 *Verifique o formato da mensagem:*\n\n` +
                `✅ *Formatos válidos:*\n` +
                `• \`23 50 pão\` (23,50 + descrição)\n` +
                `• \`100 35 teste\` (100,35 + descrição)\n` +
                `• \`pão de forma 12 23\` (descrição + 12,23)\n` +
                `• \`21 doce\` (21,00 + descrição)\n\n` +
                `📝 *Dica:* Use espaços entre valores e descrição\n` +
                `💡 *Exemplo:* "Ontem 31 00 panificadora"`;
            
            await client.sendText(message.from, errorMsg);
        }

    } catch (error) {
        console.error('Erro ao processar mensagem:', error);
        console.error('URL da API:', `${GO_API_URL}/api/processar`);
        console.error('Dados enviados:', {
            mensagem: message.body,
            autor: message.author || message.from,
            data: new Date().toLocaleDateString('pt-BR')
        });
        console.error('Detalhes do erro:', {
            message: error.message,
            code: error.code,
            status: error.response?.status,
            statusText: error.response?.statusText,
            data: error.response?.data
        });
        
        // Mensagem de erro para o usuário
        const errorMsg = `🚨 *Erro interno do sistema*\n\n` +
            `😔 Não foi possível processar sua mensagem no momento.\n\n` +
            `🔄 *Tente novamente em alguns instantes*\n` +
            `📞 Se persistir, entre em contato com o suporte.`;
        
        try {
            await client.sendText(message.from, errorMsg);
        } catch (sendError) {
            console.error('Erro ao enviar mensagem de erro:', sendError);
        }
        
        await sendLog('error', 'Erro ao processar mensagem', {
            error: error.message,
            message: message
        });
    }
}

async function processImageMessage(message) {
    try {
        console.log('Processando imagem recebida...', {
            type: message.type,
            mimetype: message.mimetype,
            mediaDataMimetype: message.mediaData?.mimetype
        });
        
        let buffer;
        let mimetype = message.mimetype || message.mediaData?.mimetype || 'image/jpeg';
        
        try {
            buffer = await client.decryptFile(message);
        } catch (decryptError) {
            console.log('Erro no decryptFile, tentando downloadFile...', decryptError.message);
            try {
                buffer = await client.downloadFile(message.id);
            } catch (downloadError) {
                console.log('Erro no downloadFile, tentando getMedia...', downloadError.message);
                const media = await client.getMedia(message);
                buffer = Buffer.from(media.data, 'base64');
            }
        }
        
        if (!buffer) {
            throw new Error('Não foi possível baixar a imagem');
        }
        
        // Gerar nome único para o arquivo
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const extension = mimetype.split('/')[1] || 'jpg';
        const filename = `img_${timestamp}_${message.from.replace('@c.us', '')}.${extension}`;
        const filepath = `../arquivos/${filename}`;
        
        // Salvar arquivo
        const fs = require('fs');
        const path = require('path');
        
        const fullPath = path.resolve(__dirname, filepath);
        
        // Criar diretório se não existir
        const dir = path.dirname(fullPath);
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, { recursive: true });
        }
        
        fs.writeFileSync(fullPath, buffer);
        
        console.log(`Imagem salva: ${fullPath} (${buffer.length} bytes)`);
        
        // Enviar confirmação
        const confirmMsg = `📸 *Imagem recebida e salva!*\n\n` +
            `📁 *Arquivo:* ${filename}\n` +
            `📊 *Tamanho:* ${(buffer.length / 1024).toFixed(1)} KB\n` +
            `📅 *Data:* ${new Date().toLocaleString('pt-BR')}\n\n` +
            `✅ *Imagem armazenada com segurança!*`;
        
        await client.sendText(message.from, confirmMsg);
        
        // Log do evento
        await sendLog('info', 'Imagem salva', {
            filename: filename,
            from: message.from,
            timestamp: timestamp,
            size: buffer.length
        });
        
    } catch (error) {
        console.error('Erro ao processar imagem:', error);
        
        const errorMsg = `❌ *Erro ao salvar imagem*\n\n` +
            `😔 Não foi possível processar a imagem enviada.\n\n` +
            `🔄 Tente enviar novamente.\n\n` +
            `🔧 *Detalhes:* ${error.message}`;
        
        try {
            await client.sendText(message.from, errorMsg);
        } catch (sendError) {
            console.error('Erro ao enviar mensagem de erro de imagem:', sendError);
        }
    }
}

// Inicializar o bot
async function start() {
    try {
        console.log('Iniciando VenomBot...');
        await sendLog('info', 'Iniciando VenomBot...');

        client = await venom.create({
            session: BOT_SESSION,
            multidevice: true,
            headless: "new", // Usar novo headless
            devtools: false,
            debug: false,
            logQR: true, // Habilitar log do QR 
            disableSpins: true, // Reduzir animações
            disableWelcome: true, // Pular welcome
            updatesLog: false, // Reduzir logs
            browserArgs: [
                '--no-sandbox',
                '--disable-setuid-sandbox',
                '--disable-dev-shm-usage',
                '--disable-accelerated-2d-canvas',
                '--no-first-run',
                '--no-zygote',
                '--disable-gpu',
                '--disable-extensions',
                '--disable-background-timer-throttling',
                '--disable-backgrounding-occluded-windows',
                '--disable-renderer-backgrounding'
            ],
            autoClose: 120000, // 2 minutos para escanear
            createPathFileToken: true, // Salvar token para reconexão
        },
        // Callback para QR Code
        (base64Qr, asciiQR, attempts, urlCode) => {
            console.log('QR Code recebido (tentativa', attempts, ')');
            console.log(asciiQR); // Mostrar no console também
            
            // Enviar QR Code para a API Go
            sendQRCodeToAPI(base64Qr, attempts, urlCode);
        },
        // Callback para status
        (statusSession, session) => {
            console.log('Status da sessão:', statusSession, 'Sessão:', session);
            sendLog('info', 'Status da sessão alterado', { status: statusSession, session });
            
            if (statusSession === 'qrReadSuccess') {
                sendLog('info', 'QR Code lido com sucesso! WhatsApp conectado.');
                updateQRCodeStatus('connected');
            } else if (statusSession === 'qrReadFail') {
                sendLog('error', 'Falha na leitura do QR Code');
                updateQRCodeStatus('failed');
            } else if (statusSession === 'autocloseCalled') {
                sendLog('warning', 'Sessão fechada automaticamente');
                updateQRCodeStatus('expired');
            }
        });

        console.log('VenomBot iniciado com sucesso!');
        await sendLog('info', 'VenomBot iniciado com sucesso');

        // Configurar handlers de eventos
        client.onMessage(async (message) => {
            // Processar mensagens individuais e de grupos
            if (message.from !== 'status@broadcast') {
                await processMessage(message);
            }
        });

        client.onAck(async (ackEvent) => {
            console.log('Status da mensagem:', ackEvent);
        });

        client.onStateChange((state) => {
            console.log('Estado alterado:', state);
            sendLog('info', 'Estado do WhatsApp alterado', { state });
        });

        // Registrar bot na API Go
        try {
            await axios.post(`${GO_API_URL}/api/bot/register`, {
                session: BOT_SESSION,
                status: 'connected',
                timestamp: new Date().toISOString()
            });
        } catch (error) {
            console.error('Erro ao registrar bot na API:', error.message);
        }

        // Handler para receber comandos da API Go
        setInterval(async () => {
            try {
                const response = await axios.get(`${GO_API_URL}/api/bot/commands/${BOT_SESSION}`);
                if (response.data && response.data.commands) {
                    for (const command of response.data.commands) {
                        await executeCommand(command);
                    }
                }
            } catch (error) {
                // Silenciar erro se não houver comandos pendentes
                if (error.response?.status !== 404) {
                    console.error('Erro ao buscar comandos:', error.message);
                }
            }
        }, 5000); // Verifica a cada 5 segundos

    } catch (error) {
        console.error('Erro ao iniciar VenomBot:', error);
        await sendLog('error', 'Erro ao iniciar VenomBot', { error: error.message });
        process.exit(1);
    }
}

// Função para executar comandos vindos da API Go
async function executeCommand(command) {
    try {
        console.log('Executando comando:', command);

        switch (command.type) {
            case 'send_message':
                await client.sendText(command.to, command.message);
                break;
            
            case 'send_media':
                if (command.mediaType === 'image') {
                    await client.sendImage(command.to, command.mediaUrl, command.filename, command.caption);
                } else if (command.mediaType === 'document') {
                    await client.sendFile(command.to, command.mediaUrl, command.filename, command.caption);
                }
                break;
            
            case 'get_contacts':
                const contacts = await client.getAllContacts();
                await axios.post(`${GO_API_URL}/api/bot/contacts`, {
                    session: BOT_SESSION,
                    contacts: contacts
                });
                break;
            
            case 'get_chats':
                const chats = await client.getAllChats();
                await axios.post(`${GO_API_URL}/api/bot/chats`, {
                    session: BOT_SESSION,
                    chats: chats
                });
                break;
        }

        // Marcar comando como executado
        await axios.delete(`${GO_API_URL}/api/bot/commands/${command.id}`);

    } catch (error) {
        console.error('Erro ao executar comando:', error);
        await sendLog('error', 'Erro ao executar comando', {
            error: error.message,
            command: command
        });
    }
}

// Graceful shutdown
process.on('SIGINT', async () => {
    console.log('Encerrando VenomBot...');
    await sendLog('info', 'Encerrando VenomBot...');
    
    if (client) {
        await client.close();
    }
    
    try {
        await axios.post(`${GO_API_URL}/api/bot/unregister`, {
            session: BOT_SESSION,
            timestamp: new Date().toISOString()
        });
    } catch (error) {
        console.error('Erro ao desregistrar bot:', error.message);
    }
    
    process.exit(0);
});

process.on('SIGTERM', async () => {
    console.log('Recebido SIGTERM, encerrando...');
    await sendLog('info', 'Recebido SIGTERM, encerrando...');
    process.exit(0);
});

// Iniciar o bot
start().catch(error => {
    console.error('Erro fatal:', error);
    process.exit(1);
});
