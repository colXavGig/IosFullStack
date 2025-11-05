
import Foundation

struct Message: Codable, Identifiable {
    let ID: Int
    let ConversationId: Int
    let SenderID: Int
    let Sender: User
    let Content: String
}
