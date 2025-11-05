
import SwiftUI

struct ConversationListView: View {
    @State private var conversations = [Conversation]()
    @State private var isShowingNewConversationView = false

    var body: some View {
        NavigationView {
            List(conversations) { conversation in
                NavigationLink(destination: ConversationDetailView(conversation: conversation)) {
                    Text(conversation.Members.first?.Username ?? "Unknown")
                }
            }
            .navigationTitle("Conversations")
            .toolbar {
                Button(action: {
                    isShowingNewConversationView = true
                }) {
                    Image(systemName: "plus")
                }
            }
            .sheet(isPresented: $isShowingNewConversationView) {
                NewConversationView()
            }
            .onAppear {
                APIService.shared.getConversations { result in
                    switch result {
                    case .success(let conversations):
                        self.conversations = conversations
                    case .failure(let error):
                        print(error.localizedDescription)
                    }
                }
            }
        }
    }
}

struct ConversationListView_Previews: PreviewProvider {
    static var previews: some View {
        ConversationListView()
    }
}
