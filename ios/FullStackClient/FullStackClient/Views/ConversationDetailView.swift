
import SwiftUI

struct ConversationDetailView: View {
    let conversation: Conversation
    @State private var messages: [Message] = []
    @State private var newMessageContent = ""

    var body: some View {
        VStack {
            List(messages) { message in
                VStack(alignment: .leading) {
                    Text(message.Sender.Username)
                        .font(.headline)
                    Text(message.Content)
                }
            }

            HStack {
                TextField("New message", text: $newMessageContent)
                    .textFieldStyle(RoundedBorderTextFieldStyle())

                Button(action: {
                    APIService.shared.sendMessage(conversationID: conversation.ID, content: newMessageContent) { result in
                        switch result {
                        case .success:
                            // Refresh messages
                            fetchMessages()
                            newMessageContent = ""
                        case .failure(let error):
                            print(error.localizedDescription)
                        }
                    }
                }) {
                    Text("Send")
                }
            }
            .padding()
        }
        .navigationTitle(conversation.Members.first?.Username ?? "Unknown")
        .onAppear {
            fetchMessages()
        }
    }

    func fetchMessages() {
        APIService.shared.getMessages(conversationID: conversation.ID) { result in
            switch result {
            case .success(let messages):
                self.messages = messages
            case .failure(let error):
                print(error.localizedDescription)
            }
        }
    }
}

struct ConversationDetailView_Previews: PreviewProvider {
    static var previews: some View {
        ConversationDetailView(conversation: Conversation(ID: 1, Members: [], Messages: []))
    }
}
