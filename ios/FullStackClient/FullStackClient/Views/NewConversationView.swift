
import SwiftUI

struct NewConversationView: View {
    @Environment(\.presentationMode) var presentationMode
    @State private var searchText = ""
    @State private var searchResults = [User]()

    var body: some View {
        NavigationView {
            VStack {
                TextField("Search for users", text: $searchText)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                    .padding()
                    .onChange(of: searchText) { newValue in
                        APIService.shared.searchUsers(query: newValue) { result in
                            switch result {
                            case .success(let users):
                                self.searchResults = users
                            case .failure(let error):
                                print(error.localizedDescription)
                            }
                        }
                    }

                List(searchResults, id: \.Username) { user in
                    Button(action: {
                        APIService.shared.createConversation(with: user.ID) { result in
                            switch result {
                            case .success:
                                presentationMode.wrappedValue.dismiss()
                            case .failure(let error):
                                print(error.localizedDescription)
                            }
                        }
                    }) {
                        Text(user.Username)
                    }
                }
            }
            .navigationTitle("New Conversation")
            .toolbar {
                Button("Cancel") {
                    presentationMode.wrappedValue.dismiss()
                }
            }
        }
    }
}

struct NewConversationView_Previews: PreviewProvider {
    static var previews: some View {
        NewConversationView()
    }
}
